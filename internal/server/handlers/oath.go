package handlers

import (
	"github.com/julienschmidt/httprouter"
	"Turn_on_PC/pkg/logging"
	"net/http"
	"encoding/json"
	"Turn_on_PC/internal/DTO"
	"Turn_on_PC/internal/server/servis"
	"Turn_on_PC/internal/server/DB"
	"Turn_on_PC/internal/server/middleware"
	"github.com/go-playground/validator/v10"
	"Turn_on_PC/internal/server/apperror"
)

const (
	oathUrl = "/oauth"
)

type handler struct {
	logger   *logging.Logger
	db       DB.DB
	validate *validator.Validate
}

func NewHandler(logger *logging.Logger, db DB.DB) Handler {
	return &handler{
		logger:   logger,
		db:       db,
		validate: validator.New(),
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPut, oathUrl, middleware.Middleware(h.SingUP))
	router.HandlerFunc(http.MethodPost, oathUrl, middleware.Middleware(h.SingIn))
}

func (h *handler) SingUP(w http.ResponseWriter, r *http.Request) error {

	user := new(DTO.CreateUser)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.Decode(&user)
	err := h.validate.Struct(user)
	if err != nil {
		return apperror.BadRequest
	}
	_, err = servis.Register(h.db, user) //todo понять нужно ли возвращать юзера или нет
	if err != nil {
		return err
	}
	w.WriteHeader(201)
	return nil
}

func (h *handler) SingIn(w http.ResponseWriter, r *http.Request) error {
	user := new(DTO.UserSingIn)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.Decode(&user)
	err := h.validate.Struct(user)
	if err != nil {
		return apperror.BadRequest
	}
	token, err := servis.SingIn(h.db, user.Login, user.Password, user.Scope)
	if err != nil {
		return err
	}
	w.WriteHeader(200)
	w.Write([]byte(token))
	return nil
}
