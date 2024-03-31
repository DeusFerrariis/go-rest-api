package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
)

type User struct {
	Name string
}

type Post struct {
	author  string
	content string
}

type Middleware func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func WithUsers(users []User) func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "users", &users)
		next(rw, r.WithContext(ctx))
	}
}

type CreateUserRequest struct {
	Username string `json:"username"`
}

func HandleCreateUser(rw http.ResponseWriter, r *http.Request) {
	users, ok := r.Context().Value("users").(*[]User)
	log.Debug("users", users)
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var userReq CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		log.Error("recieved bad request", "err", err)
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "invalid payload")
		return
	}

	for _, u := range *users {
		if u.Name == userReq.Username {
			rw.WriteHeader(http.StatusBadRequest)
			log.Debug("request for existing user", "username", userReq.Username)
			io.WriteString(rw, "user already exists")
			return
		}
	}

	*users = append(*users, User{Name: userReq.Username})
	rw.WriteHeader(http.StatusAccepted)
}
