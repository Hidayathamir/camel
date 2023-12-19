package service

import (
	"context"
	"sync"

	"github.com/Hidayathamir/camel/dto"
	"github.com/Hidayathamir/camel/repository"
)

type ChatService interface {
	SendChatRequest(ctx context.Context, payload dto.ReqStreamChat, chMessageContent chan string, wg *sync.WaitGroup)
}

type chatService struct {
	chatRepository repository.ChatRepository
}

func NewChatService(chatRepository repository.ChatRepository) ChatService {
	return &chatService{chatRepository: chatRepository}
}

func (c *chatService) SendChatRequest(ctx context.Context, payload dto.ReqStreamChat, chMessageContent chan string, wg *sync.WaitGroup) {
	c.chatRepository.SendChatRequest(ctx, payload, chMessageContent, wg)
}
