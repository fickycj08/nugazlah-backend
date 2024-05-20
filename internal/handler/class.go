package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/vandenbill/nugazlah-backend/internal/dto"
	"github.com/vandenbill/nugazlah-backend/internal/service"
	response "github.com/vandenbill/nugazlah-backend/pkg/resp"
)

type classHandler struct {
	classSvc *service.ClassService
}

func newClassHandler(classSvc *service.ClassService) *classHandler {
	return &classHandler{classSvc}
}

func (h *classHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqCreateClass

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.ResponseErrWithCode("failed to parse request body", http.StatusBadRequest, w)
		return
	}

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.classSvc.Create(r.Context(), req, token.Subject())
	if err != nil {
		response.ResponseErr(err, w)
		return
	}

	response.ResponseSuccess("success create class", nil, w)
}

func (h *classHandler) GetMyClasses(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.ResponseErrWithCode("failed to get token from request", http.StatusBadRequest, w)
		return
	}

	data, err := h.classSvc.GetMyClasses(r.Context(), token.Subject())
	if err != nil {
		response.ResponseErr(err, w)
		return
	}

	response.ResponseSuccess("success get user class", data, w)
}

func (h *classHandler) JoinClass(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.ResponseErrWithCode("failed to get token from request", http.StatusBadRequest, w)
		return
	}

	classCode := chi.URLParam(r, "classCode")
	err = h.classSvc.JoinClass(r.Context(), token.Subject(), classCode)
	if err != nil {
		response.ResponseErr(err, w)
		return
	}

	response.ResponseSuccess("success join class", nil, w)
}
