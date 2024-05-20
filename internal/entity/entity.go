package entity

import "time"

const (
	TASK_TYPE_PROPOSAL = "Proposal"
	TASK_TYPE_QUIZ     = "Quiz"
	TASK_TYPE_ESSAY    = "Essay"
	TASK_TYPE_RESPONSE = "Response"
	TASK_TYPE_PROJECT  = "Project"
)

type User struct {
	ID       string `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Class struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Lecturer    string `json:"lecturer"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Code        string `json:"code"`
	UserID      string `json:"user_id"`
}

type Task struct {
	ID      string `json:"id"`
	ClassID string `json:"class_id"`

	Title       string    `json:"title"`
	Description string    `json:"description"`
	Detail      string    `json:"detail"`
	Submission  string    `json:"submission"`
	TaskType    string    `json:"task_type"`
	Deadline    time.Time `json:"deadline"`
}

type TaskWithStatus struct {
	ID      string `json:"id"`
	ClassID string `json:"class_id"`

	Title       string    `json:"title"`
	Description string    `json:"description"`
	Detail      string    `json:"detail"`
	Submission  string    `json:"submission"`
	TaskType    string    `json:"task_type"`
	Deadline    time.Time `json:"deadline"`
	Status      bool      `json:"status"`
}

type UserClass struct {
	UserID  string `json:"user_id"`
	ClassID string `json:"class_id"`
}

type UserTask struct {
	UserID string `json:"user_id"`
	TaskID string `json:"task_id"`
	Status bool   `json:"status"`
}
