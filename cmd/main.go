package main

import (
	"context"
	"os"

	"github.com/Hidayathamir/camel/controller"
	"github.com/Hidayathamir/camel/repository"
	"github.com/Hidayathamir/camel/service"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if os.Getenv("mode") == "debug" {
		logrus.SetReportCaller(true)
	}

	chatController := initChatController()
	chatController.SendChatRequest(context.Background())
}

func initChatController() controller.ChatController {
	chatRepository := repository.NewChatRepository()
	chatService := service.NewChatService(chatRepository)
	return controller.NewChatController(chatService)
}
