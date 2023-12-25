package camel

import "path/filepath"

var (
	CamelDir     = "camel_data"
	HistoryFile  = filepath.Join(CamelDir, "history.json")
	QuestionFile = filepath.Join(CamelDir, "question.md")
	AnswerFile   = filepath.Join(CamelDir, "answer.md")
	ImagesDir    = filepath.Join(CamelDir, "images")
)
