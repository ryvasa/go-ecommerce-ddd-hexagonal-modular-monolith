package event

type TaskCreated struct {
	TaskID string
	UserID string
}

func (TaskCreated) Name() string {
	return "task.created"
}
