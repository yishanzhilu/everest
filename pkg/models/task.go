package models

// TaskStatus is the enum of status for Goal, Mission and Task
type TaskStatus uint8

const (
	// TaskTodo .
	TaskTodo TaskStatus = iota + 1
	// TaskDone .
	TaskDone
)

// TaskModel is user created scalar job.
// If status is todo, it's a todo. Todo is editable.
// IF status is done, it's a record. Record is not editable, since it will send to user's follower's timeline
type TaskModel struct {
	BaseModel
	Content string
	Review  string
	Status  TaskStatus
}
