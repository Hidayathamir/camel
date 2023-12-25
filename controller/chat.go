package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Hidayathamir/camel"
	"github.com/Hidayathamir/camel/dto"
	"github.com/Hidayathamir/camel/pkg/erra"
	"github.com/Hidayathamir/camel/service"
	"github.com/sirupsen/logrus"
)

type ChatController interface {
	SendChatRequest(ctx context.Context, model camel.Model)
}

type chatController struct {
	chatService service.ChatService
}

func NewChatController(chatService service.ChatService) ChatController {
	return &chatController{
		chatService: chatService,
	}
}

// SendChatRequest initiates a chat request handling process. It prepares the
// necessary payload, retrieves history and user questions from files, sends the
// chat request to the chat service, and manages the response streaming process.
// The incoming messages are gathered and appended to the existing history. The
// response messages are written to an answer file and updated in the history
// file for future reference.
//
// This function orchestrates the interactions between layers, considering the
// complexity of streaming responses and managing channel communication.
//
// Note: Although this function incorporates some logic that might ideally
// belong to the service layer, it accommodates complexities arising from
// response streaming and channel usage. Refactoring to adhere strictly to clean
// architecture might result in increased complexity or decreased
// maintainability due to the nature of response streaming and channel handling.
func (c *chatController) SendChatRequest(ctx context.Context, model camel.Model) {
	// this function is kinda violate clean arch, this layers should only think
	// about gathering input and then pass it into service. However since we are
	// streaming the response from llama2 and using channel to stream the
	// response, it's too complex to implement it in layer service and then
	// pass response to controller.

	payload := dto.ReqStreamChat{Model: model}

	history, err := getHistoryFromFile()
	if err != nil {
		if !errors.Is(err, camel.ErrFailedParseHistoryFile) {
			err = erra.Wrapf(err, "error get history from file")
			logrus.Error(err.Error())
			return
		}
		logrus.Warn(err.Error())
	}
	payload.Messages = history.Chat

	userQuestion, err := getUserQuestionFromFile()
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			err = erra.Wrapf(err, "error get user question from file")
			logrus.Error(err.Error())
			return
		}

		if err = createQuestionFileTemplate(); err != nil {
			err = erra.Wrapf(err, "error create question file template")
			logrus.Error(err.Error())
			return
		}

		path := filepath.Join(camel.CamelDir, camel.QuestionFile)
		logrus.Infof("please write your question in `%s` file", path)
		return
	}
	payload.Messages = append(payload.Messages, userQuestion)
	history.Chat = payload.Messages

	var wg sync.WaitGroup
	chMessageContent := make(chan string)
	defer close(chMessageContent)

	wg.Add(1)

	go c.chatService.SendChatRequest(context.Background(), payload, chMessageContent, &wg)

	var sbFullMessageResponse strings.Builder

	go func() {
		for msg := range chMessageContent {
			sbFullMessageResponse.WriteString(msg)
			fmt.Print(msg)
		}
	}()

	wg.Wait()

	fmt.Println()

	fullMsgResponse := sbFullMessageResponse.String()
	if err = writeAnswerFile(fullMsgResponse); err != nil {
		err = erra.Wrapf(err, "error write answer file")
		logrus.Error(err.Error())
		return
	}

	err = updateHistoryFile(history, fullMsgResponse)
	if err != nil {
		err = erra.Wrapf(err, "error update history file")
		logrus.Error(err.Error())
		return
	}
}

// getHistoryFromFile retrieves history data from a file located at the
// predefined camelDir. It reads the file content, parses it as JSON, and maps
// it into a History struct. If the file does not exist or encounters an error
// during reading or parsing, it returns an empty History struct and an error.
func getHistoryFromFile() (dto.History, error) {
	path := filepath.Join(camel.CamelDir, camel.HistoryFile)

	content, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return dto.History{}, erra.Wrapf(err, "error read file %s", path)
		}
		return dto.History{}, nil
	}

	var h dto.History
	if err := json.Unmarshal(content, &h); err != nil {
		err = erra.Wrapf(err, "error unmarshal content of file %s into struct History", path)
		return dto.History{}, erra.Wrap(err, camel.ErrFailedParseHistoryFile)
	}

	return h, nil
}

// getUserQuestionFromFile retrieves a user question from a file located at the
// predefined camelDir. It reads the content of the file and creates a
// dto.ReqStreamChatMessage struct with the role set as 'dto.RoleUser' and the
// content from the file.
func getUserQuestionFromFile() (dto.ReqStreamChatMessage, error) {
	path := filepath.Join(camel.CamelDir, camel.QuestionFile)

	questionContent, err := os.ReadFile(path)
	if err != nil {
		return dto.ReqStreamChatMessage{}, erra.Wrapf(err, "error read file %s", path)
	}

	userQuestion := dto.ReqStreamChatMessage{
		Role:    camel.RoleUser,
		Content: string(questionContent),
	}

	return userQuestion, nil
}

// createQuestionFileTemplate generates a template for a question file at the
// predefined camelDir location. It creates a file with the name specified by
// the questionFile constant and writes a default question string into it.
// If the file creation or writing encounters an error, it returns an error.
func createQuestionFileTemplate() error {
	if err := os.MkdirAll(camel.CamelDir, os.ModePerm); err != nil {
		return erra.Wrapf(err, "error mkdir %s", camel.CamelDir)
	}

	path := filepath.Join(camel.CamelDir, camel.QuestionFile)
	f, err := os.Create(path)
	if err != nil {
		return erra.Wrapf(err, "error create file %s", path)
	}
	defer f.Close()

	if _, err = f.WriteString("# my question\n"); err != nil {
		return erra.Wrapf(err, "error write string to file %s", path)
	}

	return nil
}

// writeAnswerFile creates a file at the predefined camelDir location and writes
// the provided answer string into it. If the file creation or writing
// encounters an error, it returns an error.
func writeAnswerFile(answer string) error {
	path := filepath.Join(camel.CamelDir, camel.AnswerFile)
	f, err := os.Create(path)
	if err != nil {
		return erra.Wrapf(err, "error create file %s", path)
	}
	defer f.Close()

	if _, err = f.WriteString(answer); err != nil {
		return erra.Wrapf(err, "error write string to file %s", path)
	}

	return nil
}

// updateHistoryFile appends a new chat message to the history and writes the
// updated history to a file located at the predefined camelDir.
func updateHistoryFile(history dto.History, content string) error {
	history.Chat = append(history.Chat, dto.ReqStreamChatMessage{
		Role:    camel.RoleAssistant,
		Content: content,
	})

	jsonByte, err := json.MarshalIndent(history, "", "    ")
	if err != nil {
		return erra.Wrapf(err, "error json marshal history")
	}

	path := filepath.Join(camel.CamelDir, camel.HistoryFile)

	if err := os.WriteFile(path, jsonByte, 0644); err != nil {
		return erra.Wrapf(err, "error write into file %s", path)
	}

	return nil
}
