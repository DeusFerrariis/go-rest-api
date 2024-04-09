package main

import (
	// "database/sql"
	"errors"
	// "github.com/charmbracelet/log"
)

type (
	User struct {
		Name string `json:"username"`
		Id   int    `json:"id"`
	}

	NewUser struct {
		Username string `json:"username"`
	}
)

// Errors

var (
	UserExistsError = errors.New("User already exists.")
)

func (us *UserStore) CreateUser(nu NewUser) (int, error) {
	// Check if exists
	var count int
	row := us.db.QueryRow("SELECT COUNT(*) as count FROM users WHERE username=?", nu.Username)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, UserExistsError
	}

	// Insert user
	// TODO: replace inline use of users with a constant table name string
	res, err := us.crud.create.Exec(nu.Username)
	if err != nil {
		return 0, err
	}

	id64, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id64), nil
}

func (us *UserStore) RetrieveUser(id int) (UserData, error) {
	var data UserData
	err := us.crud.retrieve.QueryRow(id).Scan(&data)
	if err != nil {
		return UserData{}, errors.Join(errors.New("user not found"), err)
	}
	return data, nil
}

func (us *UserStore) DeleteUser(id int) (UserData, error) {
	userData, err := us.RetrieveUser(id)
	if err != nil {
		return UserData{}, err
	}

	if _, err := us.crud.delete.Exec(id); err != nil {
		return UserData{}, errors.Join(errors.New("could not delete user"), err)
	}

	return userData, nil
}
