package dto

import "github.com/Hidayathamir/camel"

type ReqStreamChat struct {
	Model    camel.Model            `json:"model"`
	Messages []ReqStreamChatMessage `json:"messages"`
}

type ReqStreamChatMessage struct {
	Role    camel.Role
	Content string
}

type History struct {
	Chat []ReqStreamChatMessage `json:"chat"`
}

type ResStreamChat struct {
	Model     camel.Model          `json:"model"`
	CreatedAt string               `json:"created_at"`
	Message   ResStreamChatMessage `json:"message"`
	Done      bool                 `json:"done"`
}

type ResStreamChatMessage struct {
	Role    camel.Role  `json:"role"`
	Content string      `json:"content"`
	Images  interface{} `json:"images"`
}
