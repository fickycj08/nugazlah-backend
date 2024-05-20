package dto

import (
	"github.com/vandenbill/nugazlah-backend/internal/entity"
	"github.com/vandenbill/nugazlah-backend/pkg/auth"
)

type (
	ReqRegister struct {
		Email    string `json:"email" validate:"required,email"`
		FullName string `json:"fullname" validate:"required"`
		Password string `json:"password" validate:"required,min=8,max=30"`
	}
	ResRegister struct {
		UserID string `json:"user_id"`
		AccessToken string `json:"access_token"`
	}
	ReqLogin struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=30"`
	}
	ResLogin struct {
		UserID string `json:"user_id"`
		AccessToken string `json:"access_token"`
	}
)

func (d *ReqRegister) ToEntity(cryptCost int) entity.User {
	return entity.User{
		FullName: d.FullName,
		Password: auth.HashPassword(d.Password, cryptCost),
		Email:    d.Email,
	}
}
