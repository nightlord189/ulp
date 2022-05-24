package model

type TemplateMessage struct {
	Message      string
	IsError      bool
	IsAuthorized bool
	Role         string
}

func NewTmplMessage(message string, isError bool, isAuthorized bool, role string) TemplateMessage {
	return TemplateMessage{
		Message:      message,
		IsError:      isError,
		IsAuthorized: isAuthorized,
		Role:         role,
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

type TemplateEditTask struct {
	Dockerfiles []DockerfileTemplateDB
}
