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

type TemplateTasks struct {
	Tasks []TaskView
}

type TemplateAttempts struct {
	Attempts []AttemptView
}

type TemplateAttempt struct {
	IsAuthorized bool
	Role         string
	Attempt      AttemptView
}
