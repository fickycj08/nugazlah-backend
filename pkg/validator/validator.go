package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/vandenbill/nugazlah-backend/internal/entity"
)

type CredentialType string

const (
	PhoneType CredentialType = "phone"
	EmailType CredentialType = "email"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func New() *validator.Validate {
	cv := &CustomValidator{
		Validator: validator.New(),
	}

	cv.Validator.RegisterValidation("taskType", cv.taskType)

	return cv.Validator
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			err := fmt.Errorf("%s is %s", fieldError.Field(), fieldError.Tag())
			return err
		}
	}

	return nil
}

func (cv *CustomValidator) taskType(fl validator.FieldLevel) bool {
	allowedValues := []string{entity.TASK_TYPE_ESSAY, entity.TASK_TYPE_PROJECT, entity.TASK_TYPE_PROPOSAL,
		entity.TASK_TYPE_QUIZ, entity.TASK_TYPE_RESPONSE}

	value := fl.Field().String()

	for _, v := range allowedValues {
		if strings.EqualFold(value, v) {
			return true
		}
	}

	return false
}
