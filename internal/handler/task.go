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

type taskHandler struct {
	taskSvc *service.TaskService
}

func newTaskHandler(taskSvc *service.TaskService) *taskHandler {
	return &taskHandler{taskSvc}
}

func (h *taskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqCreateTask

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

	err = h.taskSvc.Create(r.Context(), req, token.Subject())
	if err != nil {
		response.ResponseErr(err, w)
		return
	}

	response.ResponseSuccess("success create task", nil, w)
}

func (h *taskHandler) GetMyTasks(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.ResponseErrWithCode("failed to get token from request", http.StatusBadRequest, w)
		return
	}

	data, err := h.taskSvc.GetMyTasks(r.Context(), token.Subject(), chi.URLParam(r, "classID"))
	if err != nil {
		response.ResponseErr(err, w)
		return
	}

	response.ResponseSuccess("success get user tasks", data, w)
}

func (h *taskHandler) GetDetailTask(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.ResponseErrWithCode("failed to get token from request", http.StatusBadRequest, w)
		return
	}

	task, err := h.taskSvc.GetDetailTask(r.Context(), token.Subject(), chi.URLParam(r, "taskID"))
	if err != nil {
		response.ResponseErr(err, w)
		return
	}

	response.ResponseSuccess("success get detail task", task, w)
}

func (h *taskHandler) MarkAsDone(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.ResponseErrWithCode("failed to get token from request", http.StatusBadRequest, w)
		return
	}

	err = h.taskSvc.SetToDone(r.Context(), token.Subject(), chi.URLParam(r, "taskID"))
	if err != nil {
		response.ResponseErr(err, w)
		return
	}

	response.ResponseSuccess("success mark as done", nil, w)
}
