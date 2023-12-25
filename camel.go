package camel

type Model string

const (
	ModelLlama27b Model = "llama2:7b"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)
