package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/vandenbill/nugazlah-backend/internal/cfg"
	"github.com/vandenbill/nugazlah-backend/internal/dto"
	"github.com/vandenbill/nugazlah-backend/internal/entity"
	"github.com/vandenbill/nugazlah-backend/internal/ierr"
	"github.com/vandenbill/nugazlah-backend/internal/repo"
)

type TaskService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newTaskService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *TaskService {
	return &TaskService{repo, validator, cfg}
}

func (u *TaskService) Create(ctx context.Context, body dto.ReqCreateTask, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}
	return u.repo.Task.Insert(ctx, sub, body.ToEntity())
}

func (u *TaskService) GetMyTasks(ctx context.Context, sub, classID string) ([]dto.ResGetMyTasks, error) {
	tasks, err := u.repo.Task.GetMyTasks(ctx, sub, classID)

	resp := make([]dto.ResGetMyTasks, 0)
	for _, v := range tasks {
		res := dto.ResGetMyTasks{}
		res.FromEntity(v)
		resp = append(resp, res)
	}
	return resp, err
}

func (u *TaskService) GetDetailTask(ctx context.Context, sub, taskID string) (entity.TaskWithStatus, error) {
	return u.repo.Task.GetTask(ctx, sub, taskID)
}

func (u *TaskService) SetToDone(ctx context.Context, sub, taskID string) error {
	return u.repo.Task.MarkTaskDone(ctx, sub, taskID)
}
