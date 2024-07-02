package main

import "database/sql"

type Store interface {

	// USERS
	CreateUser() error

	// TASKS
	CreateTask(task *Task) (*Task, error)
	GetTask(string) (*Task, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser() error {
	return nil
}

func (s *Storage) CreateTask(t *Task) (*Task, error) {
	rows, err := s.db.Exec("INSERT INTO task (name, status, projectID, assignedToID) VALUES (?, ?, ?, ?)", t.Name, t.Status, t.ProjectID, t.AssignetToID)

	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	t.ID = id
	return t, nil
}

func (s *Storage) GetTask(id string) (*Task, error) {
	var t Task

	err := s.db.QueryRow("SELECT id, name, status, projectID, assignedToID, createdAt FROM task where id = ?", id).Scan(
		&t.ID,
		&t.Name,
		&t.Status,
		&t.ProjectID,
		&t.AssignetToID,
		&t.CreatedAt,
	)

	return &t, err
}
