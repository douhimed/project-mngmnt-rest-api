package main

import "database/sql"

type Store interface {

	// USERS
	CreateUser(*User) (*User, error)
	GetUserByID(string) (*User, error)

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

func (s *Storage) CreateUser(u *User) (*User, error) {
	rows, err := s.db.Exec("INSERT INTO user (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", u.FirstName, u.LastName, u.Email, u.Password)

	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id
	return u, nil
}

func (s *Storage) GetUserByID(id string) (*User, error) {
	var u User

	err := s.db.QueryRow("SELECT id, firstName, lastName, email, createAt FROM user where id = ?", id).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.CreateAt,
	)

	return &u, err
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
