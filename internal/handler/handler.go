package handler

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/vandenbill/nugazlah-backend/internal/cfg"
	"github.com/vandenbill/nugazlah-backend/internal/service"
)

type Handler struct {
	router  *chi.Mux
	service *service.Service
	cfg     *cfg.Cfg
}

func NewHandler(router *chi.Mux, service *service.Service, cfg *cfg.Cfg) *Handler {
	handler := &Handler{router, service, cfg}
	handler.registRoute()

	return handler
}

func (h *Handler) registRoute() {
	r := h.router
	var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(h.cfg.JWTSecret), nil, jwt.WithAcceptableSkew(30*time.Second))

	userH := newUserHandler(h.service.User)
	classH := newClassHandler(h.service.Class)
	taskH := newTaskHandler(h.service.Task)

	r.Post("/v1/auth/register", userH.Register)
	r.Post("/v1/auth/login", userH.Login)

	// protected route
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/v1/classes", classH.Create)
		r.Post("/v1/classes/{classCode}/join", classH.JoinClass)
		r.Get("/v1/classes", classH.GetMyClasses)

		r.Post("/v1/tasks", taskH.Create)
		r.Post("/v1/tasks/{taskID}/done", taskH.MarkAsDone)
		r.Get("/v1/tasks/classes/{classID}", taskH.GetMyTasks)
		r.Get("/v1/tasks/{taskID}", taskH.GetDetailTask)
	})
}
