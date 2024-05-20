package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/vandenbill/nugazlah-backend/internal/cfg"
	"github.com/vandenbill/nugazlah-backend/internal/dto"
	"github.com/vandenbill/nugazlah-backend/internal/ierr"
	"github.com/vandenbill/nugazlah-backend/internal/repo"
	"github.com/vandenbill/nugazlah-backend/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newUserService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *UserService {
	return &UserService{repo, validator, cfg}
}

func (u *UserService) Login(ctx context.Context, body dto.ReqLogin) (dto.ResLogin, error) {
	res := dto.ResLogin{}

	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	user, err := u.repo.User.FindByEmail(ctx, body.Email)
	if err != nil {
		return res, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 10, auth.JwtPayload{Sub: user.ID})
	if err != nil {
		return res, err
	}

	res.UserID = user.ID
	res.AccessToken = token

	return res, nil
}

func (u *UserService) Register(ctx context.Context, body dto.ReqRegister) (dto.ResRegister, error) {
	res := dto.ResRegister{}

	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	user, err := u.repo.User.FindByEmail(ctx, body.Email)
	if err != nil && err != ierr.ErrNotFound {
		return res, err
	}

	err = u.repo.User.Insert(ctx, body.ToEntity(10))
	if err != nil {
		return res, err
	}

	user, err = u.repo.User.FindByEmail(ctx, body.Email)
	if err != nil {
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 10, auth.JwtPayload{Sub: user.ID})
	if err != nil {
		return res, err
	}

	res.UserID = user.ID
	res.AccessToken = token
	return res, nil
}
