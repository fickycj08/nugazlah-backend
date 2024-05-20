package dto

import "github.com/vandenbill/nugazlah-backend/internal/entity"

type (
	ReqCreateClass struct {
		Name        string `json:"name" validate:"required,min=3"`
		Lecturer    string `json:"lecturer" validate:"required,min=3"`
		Description string `json:"description" validate:"required,min=10"`
		Icon        string `json:"icon" validate:"required"`
	}

	ReqJoinClass struct {
		ClassCode string `json:"class_code" validate:"required,min=6,max=6"`
	}

	ResGetMyClasses struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Lecturer    string `json:"lecturer"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Code        string `json:"code"`
		Maker       string `json:"maker"`
	}
)

func (r *ResGetMyClasses) FromEntity(e entity.Class) {
	r.ID = e.ID
	r.Name = e.Name
	r.Lecturer = e.Lecturer
	r.Description = e.Description
	r.Icon = e.Icon
	r.Code = e.Code
	r.Maker = e.UserID
}
