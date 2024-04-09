package main

import (
	"database/sql"
	"sync"

	// "github.com/charmbracelet/log"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

const (
	UserFile         string = "users.db"
	CreateUsersTable string = `
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE
		);`
)

type CRUD struct {
	create, retrieve, delete *sql.Stmt
}

type UserStore struct {
	mu   sync.Mutex
	db   *sql.DB
	crud *CRUD
}

func NewUserStore() (*UserStore, error) {
	log.Info("initializing new UserStore", "db_path", UserFile)
	db, err := sql.Open("sqlite3", UserFile)
	if err != nil {
		log.Error("could not initialize sqlite connection", "err", err)
		return nil, err
	}
	if _, err := db.Exec(CreateUsersTable); err != nil {
		log.Error("Error configuring database", "err", err)
		return nil, err
	}
	crud, err := NewUserStoreCRUD(db)
	if err != nil {
		return nil, err
	}

	log.Info("Database connected & configured for users")
	return &UserStore{
		db:   db,
		crud: crud,
	}, nil
}

func NewUserStoreCRUD(db *sql.DB) (*CRUD, error) {
	log.Debug("Adding prepared crud sql statements to user store db")
	create, err := db.Prepare("INSERT INTO users VALUES(NULL, ?);")
	if err != nil {
		return nil, err
	}
	retrieve, err := db.Prepare("SELECT * FROM users WHERE id=?;")
	if err != nil {
		return nil, err
	}
	delete, err := db.Prepare("SELECT * FROM users WHERE id=?;")
	if err != nil {
		return nil, err
	}

	return &CRUD{
		create:   create,
		retrieve: retrieve,
		delete:   delete,
	}, nil
}
