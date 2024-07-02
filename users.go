package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type UserService struct {
	store Store
}

func NewUSerService(s Store) *UserService {
	return &UserService{
		store: s,
	}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.handleUserRegister).Methods("POST")
}

func (s *UserService) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload *User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid requedt payload",
		})
		return
	}

	if err := validateUserPayload(payload); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	hashedPwd, err := HashPassword(payload.Password)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "error hashing password",
		})
		return
	}

	payload.Password = hashedPwd

	u, err := s.store.CreateUser(payload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "error creating user",
		})
		return
	}

	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "eror creating session",
		})
		return
	}

	WriteJSON(w, http.StatusCreated, token)
}

var errEmailRequired = errors.New("email is required")
var errFirstNameRequired = errors.New("firstname is required")
var errLastNameRequired = errors.New("lastname is required")
var errPasswordRequired = errors.New("password is required")

func validateUserPayload(u *User) error {
	if u.Email == "" {
		return errEmailRequired
	}

	if u.FirstName == "" {
		return errFirstNameRequired
	}

	if u.LastName == "" {
		return errLastNameRequired
	}

	if u.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secret := []byte(Envs.JWTSecret)
	token, err := CreateJWT(secret, id)
	if err != nil {
		return "", nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}
