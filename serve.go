package main

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/urfave/negroni"
)

func main() {
	// Routing
	// Middleware
	// State
	users := make([]User, 0)
	withUsers := WithUsers(&users)
	n := negroni.New()
	n.Use(negroni.HandlerFunc(LogMiddleware))
	n.Use(negroni.HandlerFunc(withUsers))

	r := chi.NewRouter()
	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		io.WriteString(rw, "Hello, world!")
	})
	r.Get("/user", HandleGetUser)
	r.Post("/user/new", HandleCreateUser)
	r.Post("/post/new", HandleCreatePost)
	r.Get("/user/posts", HandleGetUserPosts)
	r.Delete("/user", HandleDeleteUser)

	n.UseHandler(r)

	http.ListenAndServe(":3001", n)
}
