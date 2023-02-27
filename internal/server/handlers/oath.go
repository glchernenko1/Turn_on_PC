package handlers

import (
	"Turn_on_PC/internal/DTO"
	"Turn_on_PC/internal/server/DB"
	"Turn_on_PC/internal/server/apperror"
	"Turn_on_PC/internal/server/middleware"
	"Turn_on_PC/internal/server/servis"
	"Turn_on_PC/pkg/logging"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"net/http"
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
	h.logger.Info("start SingIn ")
	user := new(DTO.UserSingIn)
	decoder := json.NewDecoder(r.Body)
	h.logger.Info("Decode json")
	defer r.Body.Close()
	decoder.Decode(&user)
	h.logger.Info("Decode User")
	err := h.validate.Struct(user)
	h.logger.Info("Validate")
	if err != nil {
		return apperror.BadRequest
	}
	token, err := servis.SingIn(h.db, user.Login, user.Password, user.Scope)
	h.logger.Info("create token")
	if err != nil {
		return err
	}
	w.WriteHeader(200)
	w.Write([]byte(token))
	h.logger.Info("push token")
	return nil
}
