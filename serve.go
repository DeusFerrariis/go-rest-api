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
	withUsers := WithUsers(users)
	n := negroni.New()
	n.Use(negroni.HandlerFunc(LogMiddleware))
	n.Use(negroni.HandlerFunc(withUsers))

	r := chi.NewRouter()
	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		io.WriteString(rw, "Hello, world!")
	})
	r.Post("/user/new", HandleCreateUser)

	n.UseHandler(r)

	http.ListenAndServe(":3001", n)
}
