package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vandenbill/nugazlah-backend/internal/dto"
	"github.com/vandenbill/nugazlah-backend/internal/service"
	response "github.com/vandenbill/nugazlah-backend/pkg/resp"
)

type userHandler struct {
	userSvc *service.UserService
}

func newUserHandler(userSvc *service.UserService) *userHandler {
	return &userHandler{userSvc}
}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqRegister

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.ResponseErrWithCode("failed to parse request body", http.StatusBadRequest, w)
		return
	}

	res, err := h.userSvc.Register(r.Context(), req)
	if err != nil {
		response.ResponseErr(err, w)
		return
	}

	response.ResponseSuccess("User registered successfully", res, w)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqLogin

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.ResponseErrWithCode("failed to parse request body", http.StatusBadRequest, w)
		return
	}

	res, err := h.userSvc.Login(r.Context(), req)
	if err != nil {
		response.ResponseErr(err, w)
		return
	}

	response.ResponseSuccess("User logged successfully", res, w)
}
