package camel

import "errors"

var (
	ErrFailedParseHistoryFile = errors.New("error failed parse history file")
	ErrQuestionFileNotFound   = errors.New("error question file not found")
	ErrImagesDiNotFound       = errors.New("error images dir not found")
)
