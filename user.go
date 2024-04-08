package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
	resp "github.com/nicklaw5/go-respond"
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

	UserIdRequest struct {
		UserId int `json:"user_id"`
	}
)

func FromRequest[V interface{}](r *http.Request) (*V, error) {
	var body V
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

type BodyHandler[B interface{}] func(http.ResponseWriter, *http.Request, B)

func HandleWithBody[B interface{}](next BodyHandler[B]) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		body, err := FromRequest[B](r)
		if err != nil {
			resp.NewResponse(rw).BadRequest(nil)
			log.Error("recieved bad request: %s", err)
			return
		}
		next(rw, r, *body)
	}
}

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

func Message(msg string) map[string]string {
	return map[string]string{"message": msg}
}

func handleCreateUser(rw http.ResponseWriter, r *http.Request, b CreateUserRequest) {
	if b.Username == "" {
		resp.NewResponse(rw).BadRequest(Message("username can't be empty"))
		return
	}

	users, ok := r.Context().Value(ContextKey{"users"}).(*[]User)
	if !ok {
		resp.NewResponse(rw).DefaultMessage().InternalServerError(nil)
		return
	}

	for _, u := range *users {
		if u.Name == b.Username {
			resp.NewResponse(rw).BadRequest(Message("user with that username, already exists"))
			log.Debug("request for existing user", "username", b.Username)
			return
		}
	}

	id := 1
	if len(*users) != 0 {
		id = (*users)[len(*users)-1].Id + 1
	}

	*users = append(*users, User{Name: b.Username, Id: id})

	resp.NewResponse(rw).Accepted(map[string]int{"user_id": id})
}

var HandleCreateUser = HandleWithBody[CreateUserRequest](handleCreateUser)

func handleDeleteUser(rw http.ResponseWriter, r *http.Request, b UserIdRequest) {
	users, ok := r.Context().Value(ContextKey{"users"}).(*[]User)
	if !ok {
		resp.NewResponse(rw).DefaultMessage().InternalServerError(nil)
		log.Error("could not retrieve user from request context")
		return
	}

	for i, user := range *users {
		if user.Id == b.UserId {
			// All except user at position i
			*users = append((*users)[:i], (*users)[i+1:]...)
		}
	}
}

var HandleDeleteUser = HandleWithBody[UserIdRequest](handleDeleteUser)

// TODO  convert to use HandleWithBody
func HandleCreatePost(rw http.ResponseWriter, r *http.Request) {
	reqBody, err := FromRequest[CreatePostRequest](r)
	if err != nil {
		log.Error("recieved bad request", "err", err)
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "invalid payload")
		return
	}

	p := Post{
		Author:  reqBody.Username,
		Content: reqBody.Content,
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
				log.Errorf("could not serialize user: %s", err.Error())
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
