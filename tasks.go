package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type TaskService struct {
	store Store
}

func NewTaskService(s Store) *TaskService {
	return &TaskService{
		store: s,
	}
}

func (ts *TaskService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", ts.handleCreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", ts.handleGetTask).Methods("GET")
}

func (ts *TaskService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "invalide request payload"})
		return
	}

	defer r.Body.Close()

	var task *Task

	err = json.Unmarshal(body, task)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalide request payload"})
		return
	}

	if err := validateTask(task); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	t, err := ts.store.CreateTask(task)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating task"})
		return
	}

	WriteJSON(w, http.StatusCreated, t)
}

func (ts *TaskService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	t, err := ts.store.GetTask(id)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: "task not found"})
		return
	}

	WriteJSON(w, http.StatusOK, t)
}

var errNameRequired = errors.New("name is required")
var errProjectIdRequired = errors.New("project id is required")
var errUserIdRequired = errors.New("user id is required")

func validateTask(task *Task) error {

	if task.Name == "" {
		return errNameRequired
	}

	if task.ProjectID <= 0 {
		return errProjectIdRequired
	}

	if task.AssignetToID <= 0 {
		return errUserIdRequired
	}

	return nil
}
