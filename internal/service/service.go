package service

import (
	"github.com/go-playground/validator/v10"

	"github.com/vandenbill/nugazlah-backend/internal/cfg"
	"github.com/vandenbill/nugazlah-backend/internal/repo"
)

type Service struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg

	User   *UserService
	Class   *ClassService
	Task   *TaskService
}

func NewService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *Service {
	service := Service{}
	service.repo = repo
	service.validator = validator
	service.cfg = cfg

	service.User = newUserService(repo, validator, cfg)
	service.Class = newClassService(repo, validator, cfg)
	service.Task = newTaskService(repo, validator, cfg)

	return &service
}
