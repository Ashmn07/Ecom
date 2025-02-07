package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ashmn07/Ecom/service/auth"
	"github.com/Ashmn07/Ecom/types"
	"github.com/Ashmn07/Ecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
}
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//Get JSON Payload

	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	//validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	//Check if user Exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)

	if err != nil {
		log.Print("H6")
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//Else Create new user
	err = h.store.CreateUser(types.User{
		Firstname: payload.Firstname,
		Lastname:  payload.Lastname,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		log.Fatal(err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
