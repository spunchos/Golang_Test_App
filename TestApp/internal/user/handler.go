package user

import (
	"TestApp/internal/apperror"
	"TestApp/internal/handlers"
	"TestApp/pkg/logging"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	usersURL = "/users"
	userURL  = "/users/:id"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

//TODO make a service interface

func NewHandler(service *Service, logger *logging.Logger) handlers.Handler {
	return &handler{service: service, logger: logger}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetUsers))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUser))
	router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.DeleteUser))
}

//TODO add valiadation

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.service.FindAll(r.Context())
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(users)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return apperror.NewAppError(err, "Invalid input", "Invalid input", "US-000001")
	}
	user, err := h.service.Create(r.Context(), &u)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(user)
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	user, err := h.service.FindOne(r.Context(), id)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(user)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return apperror.NewAppError(err, "Invalid input", "Invalid input", "US-000002")
	}
	u.ID = id
	if err := h.service.Update(r.Context(), u); err != nil {
		return err
	}
	return nil
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if err := h.service.Delete(r.Context(), id); err != nil {
		return err
	}
	return nil
}
