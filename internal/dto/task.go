package dto

import (
	"time"

	"github.com/vandenbill/nugazlah-backend/internal/entity"
)

type (
	ReqCreateTask struct {
		ClassID     string    `json:"class_id" validate:"uuid"`
		Title       string    `json:"title" validate:"required,min=3"`
		Detail      string    `json:"detail" validate:"required,min=3"`
		Description string    `json:"description" validate:"required,min=10"`
		Submission  string    `json:"submission" validate:"required,min=3"`
		TaskType    string    `json:"task_type" validate:"required,taskType"`
		Deadline    time.Time `json:"deadline" validate:"required"`
	}
	ResGetMyTasks struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		TaskType    string    `json:"task_type"`
		Deadline    time.Time `json:"deadline"`
		Description string    `json:"description"`
		IsDone      bool      `json:"is_done"`
	}
)

func (r *ResGetMyTasks) FromEntity(e entity.TaskWithStatus) {
	r.ID = e.ID
	r.Title = e.Title
	r.TaskType = e.TaskType
	r.Deadline = e.Deadline
	r.IsDone = e.Status
	r.Description = e.Description
}

func (r *ReqCreateTask) ToEntity() entity.Task {
	return entity.Task{
		ClassID:     r.ClassID,
		Title:       r.Title,
		Detail:      r.Detail,
		Description: r.Description,
		Submission:  r.Submission,
		TaskType:    r.TaskType,
		Deadline:    r.Deadline,
	}
}
