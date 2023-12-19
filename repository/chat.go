package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Hidayathamir/camel/dto"
	"github.com/Hidayathamir/camel/pkg/erra"
	"github.com/sirupsen/logrus"
)

type ChatRepository interface {
	SendChatRequest(ctx context.Context, payload dto.ReqStreamChat, chMessageContent chan string, wg *sync.WaitGroup)
}

type chatRepository struct{}

func NewChatRepository() ChatRepository {
	return &chatRepository{}
}

func (*chatRepository) SendChatRequest(ctx context.Context, payload dto.ReqStreamChat, chMessageContent chan string, wg *sync.WaitGroup) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		err = erra.Wrapf(err, "error marshal payload")
		logrus.WithFields(logrus.Fields{"payload": payload}).Error(err.Error())
		wg.Done()
		return
	}

	url := "http://localhost:11434/api/chat"
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		err = erra.Wrapf(err, "error http post request")
		logrus.WithFields(logrus.Fields{
			"url":     url,
			"payload": payload,
		}).Error(err.Error())
		wg.Done()
		return
	}
	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)
	for {
		var r dto.ResStreamChat
		if err := dec.Decode(&r); err != nil {
			break
		}
		chMessageContent <- r.Message.Content
	}

	wg.Done()
}
