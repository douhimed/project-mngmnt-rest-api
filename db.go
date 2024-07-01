package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(cfg mysql.Config) *MySQLStorage {
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to mysql OK")

	return &MySQLStorage{
		db: db,
	}
}

func (s *MySQLStorage) Init() (*sql.DB, error) {

	if err := s.createProjectsTable(); err != nil {
		return nil, err
	}

	if err := s.createUSersTable(); err != nil {
		return nil, err
	}

	if err := s.createTasksTable(); err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *MySQLStorage) createUSersTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			firstName VARCHAR(255) NOT NULL,
			lastName VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			createAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (id),
			UNIQUE KEY (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8`)

	return err
}

func (s *MySQLStorage) createProjectsTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS project (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			createAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8`)

	return err
}

func (s *MySQLStorage) createTasksTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS task (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			createAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			status ENUM('TODO', 'IN_PROGRESS', 'IN_TESTING', 'DONE') NOT NULL DEFAULT 'TODO',
			projectID INT UNSIGNED NOT NULL,
			assignedToID INT UNSIGNED NOT NULL,

			PRIMARY KEY (id),
			FOREIGN KEY (projectID) REFERENCES project(id),
			FOREIGN KEY (assignedToID) REFERENCES user(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8`)

	return err
}
