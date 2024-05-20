package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/vandenbill/nugazlah-backend/internal/cfg"
	"github.com/vandenbill/nugazlah-backend/internal/dto"
	"github.com/vandenbill/nugazlah-backend/internal/entity"
	"github.com/vandenbill/nugazlah-backend/internal/ierr"
	"github.com/vandenbill/nugazlah-backend/internal/repo"
	"github.com/vandenbill/nugazlah-backend/pkg/random"
)

type ClassService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newClassService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *ClassService {
	return &ClassService{repo, validator, cfg}
}

func (u *ClassService) Create(ctx context.Context, body dto.ReqCreateClass, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	err = u.repo.Class.Insert(ctx, entity.Class{
		Name:        body.Name,
		Lecturer:    body.Lecturer,
		Description: body.Description,
		Icon:        body.Icon,
		Code:        random.GenerateRandomCapitalString(6),
		UserID:      sub,
	})

	return err
}

func (u *ClassService) GetMyClasses(ctx context.Context, sub string) ([]dto.ResGetMyClasses, error) {
	resp := make([]dto.ResGetMyClasses, 0)
	classes, err := u.repo.Class.GetMyClasses(ctx, sub)
	if err != nil {
		return nil, err
	}
	for _, v := range classes {
		res := dto.ResGetMyClasses{}
		res.FromEntity(v)
		resp = append(resp, res)
	}
	return resp, err
}

func (u *ClassService) GetClass(ctx context.Context, classID string) (entity.Class, error) {
	class, err := u.repo.Class.GetClass(ctx, classID)
	return class, err
}

func (u *ClassService) JoinClass(ctx context.Context, userID, classCode string) error {
	isAlreadyJoin, err := u.repo.Class.IsAlreadyJoin(ctx, userID, classCode)
	if err != nil {
		return err
	}
	if isAlreadyJoin {
		return ierr.ErrAlreadyJoinClass
	}
	return u.repo.Class.JoinClass(ctx, userID, classCode)
}
