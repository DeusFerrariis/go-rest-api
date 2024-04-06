package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
)

type (
	User struct {
		Name string `json:"username`
		Id   int    `json:"id"`
	}

	Post struct {
		Author  string `json:"author"`
		Content string `json:"content"`
	}

	Dict map[string]interface{}

	Middleware func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

	ContextKey struct {
		string
	}

	CreateUserRequest struct {
		Username string `json:"username"`
	}

	CreatePostRequest struct {
		Username string `json:"username"`
		Content  string `json:"content"`
	}
)

// Middlewares

func WithUsers(users *[]User) func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKey{"users"}, users)
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

// Handlers

func HandleDeleteUser(rw http.ResponseWriter, r *http.Request) {
	users, ok := r.Context().Value(ContextKey{"users"}).(*[]User)

	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error("could not retrieve user from request context")
	}

	var id string
	if id = r.URL.Query().Get("id"); id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "missing name parameter in URL query")
		return
	}
	var idInt int
	if i, err := strconv.Atoi(id); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "invalid id"+id)
		return
	} else {
		idInt = i
	}

	for i, user := range *users {
		if user.Id == idInt {
			// All except user at position i
			*users = append((*users)[:i], (*users)[i+1:]...)
		}
	}
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

	var id int
	if len(*users) == 0 {
		id = 1
	} else {
		id = (*users)[len(*users)-1].Id + 1
	}

	*users = append(*users, User{Name: userReq.Username, Id: id})

	if err := json.NewEncoder(rw).Encode(map[string]int{"user_id": id}); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Errorf("Error serializing: %s", err.Error())
		return
	}
	rw.WriteHeader(http.StatusAccepted)
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
	var id string
	if id = r.URL.Query().Get("id"); id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "missing name parameter in URL query")
		return
	}
	var idInt int
	if i, err := strconv.Atoi(id); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "invalid id"+id)
		return
	} else {
		idInt = i
	}

	users, ok := r.Context().Value(ContextKey{"users"}).(*[]User)
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error("could not retrieve users from request context")
	}

	for _, user := range *users {
		if idInt == user.Id {
			if err := json.NewEncoder(rw).Encode(dataResponse(user)); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				log.Errorf("could not serialize user: ", err.Error())
			}
			return
		}
	}

	rw.WriteHeader(http.StatusNotFound)
	io.WriteString(rw, "user not found")
}

func HandleGetUserPosts(rw http.ResponseWriter, r *http.Request) {
	var name string
	if name = r.URL.Query().Get("username"); name == "" {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "missing name parameter in URL query")
		return
	}

	posts, ok := r.Context().Value(ContextKey{"posts"}).(*[]Post)
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error("could not retrieve users from request context")
		return
	}

	var userPosts []Post = make([]Post, 0)
	for _, post := range *posts {
		if post.Author == name {
			userPosts = append(userPosts, post)
		}
	}

	json.NewEncoder(rw).Encode(dataResponse(userPosts))
}

func dataResponse(v any) Dict {
	return map[string]interface{}{
		"data": v,
	}
}
