package camel

type Model string

const (
	ModelLlama27b          Model = "llama2:7b"
	ModelLlava7b           Model = "llava:7b"
	ModelCodeLlama34b      Model = "codellama:34b"
	ModelCodeLlamaInstruct Model = "codellama:instruct"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)
