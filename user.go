package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

type User struct {
	Name string
}

type Post struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type Middleware func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

type ContextKey struct {
	string
}

func WithUsers(users []User) func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKey{"users"}, &users)
		next(rw, r.WithContext(ctx))
	}
}

func WithPosts(posts []Post) func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKey{"posts"}, &posts)
		next(rw, r.WithContext(ctx))
	}
}

type CreateUserRequest struct {
	Username string `json:"username"`
}

func HandleCreateUser(rw http.ResponseWriter, r *http.Request) {
	users, ok := r.Context().Value(ContextKey{"users"}).(*[]User)
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

type CreatePostRequest struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func HandleCreatePost(rw http.ResponseWriter, r *http.Request) {
	var postReq CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&postReq); err != nil {
		log.Error("recieved bad request", "err", err)
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "invalid payload")
		return
	}

	p := Post{
		Author:  postReq.Username,
		Content: postReq.Content,
	}

	posts, ok := r.Context().Value(ContextKey{"posts"}).(*[]Post)
	log.Debug("posts", posts)
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error("could not retrieve posts from request context")
		return
	}

	*posts = append(*posts, p)
}

// func HandleCreatePost(rw http.ResponseWriter, r *http.Request) {
func HandleGetUser(rw http.ResponseWriter, r *http.Request) {
	if name := chi.URLParam(r, "username"); name == "" {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "missing name parameter in URL query")
		return
	}

	
}

func HandleGetUserPosts(rw http.ResponseWriter, r *http.Request) {
	if name := chi.URLParam(r, "username"); name == "" {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "missing name parameter in URL query")
		return
	}
}
