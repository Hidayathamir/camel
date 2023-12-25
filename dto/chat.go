package dto

import "github.com/Hidayathamir/camel"

type ReqStreamChat struct {
	Model    camel.Model            `json:"model"`
	Messages []ReqStreamChatMessage `json:"messages"`
}

type Base64String string

type ReqStreamChatMessage struct {
	Role    camel.Role     `json:"role"`
	Content string         `json:"content"`
	Images  []Base64String `json:"images"`
}

type History struct {
	Chat []ReqStreamChatMessage `json:"chat"`
}

type stringDateTime string

type ResStreamChat struct {
	Model     camel.Model          `json:"model"`
	CreatedAt stringDateTime       `json:"created_at"`
	Message   ResStreamChatMessage `json:"message"`
	Done      bool                 `json:"done"`
}

type ResStreamChatMessage struct {
	Role    camel.Role  `json:"role"`
	Content string      `json:"content"`
	Images  interface{} `json:"images"`
}
