package model

type TemplateUserInfo struct {
	IsAuthorized bool
	Username     string
	Role         string
}

type TemplateMessage struct {
	UserInfo TemplateUserInfo
	Message  string
	IsError  bool
}

func NewTmplMessage(message string, isError bool, isAuthorized bool, username, role string) TemplateMessage {
	return TemplateMessage{
		Message: message,
		IsError: isError,
		UserInfo: TemplateUserInfo{
			IsAuthorized: isAuthorized,
			Username:     username,
			Role:         role,
		},
	}
}

type TemplateTasks struct {
	UserInfo TemplateUserInfo
	Tasks    []TaskView
}

type TemplateAttempts struct {
	UserInfo TemplateUserInfo
	Attempts []AttemptView
}

type TemplateAttempt struct {
	UserInfo TemplateUserInfo
	Attempt  AttemptView
}

type TemplateEditTask struct {
	UserInfo         TemplateUserInfo
	UserID           int
	IsEdit           bool
	ID               int
	Name             string
	Description      string
	TaskType         string
	Dockerfile       string
	TestCaseType     TestCaseType
	TestCaseURL      string
	TestCaseInput    string
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
	t.TestCaseInput = task.TestcaseInput
	t.TestCaseExpected = task.TestcaseExpected
}

type TemplateUploadAttempt struct {
	UserInfo        TemplateUserInfo
	UserID          int
	TaskID          int
	TaskName        string
	TaskDescription string
	TaskType        string
}

func (t *TemplateUploadAttempt) Fill(task TaskDB) {
	t.TaskID = task.ID
	t.TaskName = task.Name
	t.TaskDescription = task.Description
	t.TaskType = string(task.Type)
}
