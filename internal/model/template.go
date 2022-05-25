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
	IsEdit           bool
	ID               int
	Name             string
	Description      string
	TaskType         string
	Dockerfile       string
	TestCaseType     string
	TestCaseURL      string
	TestCaseExpected string
	Dockerfiles      []DockerfileTemplateDB
}

func (t *TemplateEditTask) Fill(task TaskDB) {
	t.ID = task.ID
	t.Name = task.Name
	t.Description = task.Description
	t.TaskType = string(task.Type)
	t.Dockerfile = task.Dockerfile
	t.TestCaseType = task.TestcaseType
	t.TestCaseURL = task.TestcaseURL
	t.TestCaseExpected = task.TestcaseExpected
}