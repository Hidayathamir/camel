package dto

type Model string

const (
	ModelLlama27b Model = "llama2:7b"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type ReqStreamChat struct {
	Model    Model                  `json:"model"`
	Messages []ReqStreamChatMessage `json:"messages"`
}

type ReqStreamChatMessage struct {
	Role    Role
	Content string
}

type ResStreamChat struct {
	Model     Model                `json:"model"`
	CreatedAt string               `json:"created_at"`
	Message   ResStreamChatMessage `json:"message"`
	Done      bool                 `json:"done"`
}

type ResStreamChatMessage struct {
	Role    Role        `json:"role"`
	Content string      `json:"content"`
	Images  interface{} `json:"images"`
}
