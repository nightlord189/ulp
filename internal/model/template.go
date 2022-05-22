package model

type TemplateMessage struct {
	Message string
	IsError bool
}

func NewTmplMessage(message string, isError bool) TemplateMessage {
	return TemplateMessage{
		Message: message,
		IsError: isError,
	}
}
