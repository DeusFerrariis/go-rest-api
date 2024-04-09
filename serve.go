package main

import (
	"io"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/urfave/negroni"
)

func main() {
	// Routing
	// Middleware
	// State
	n := negroni.New()
	n.Use(negroni.HandlerFunc(LogMiddleware))

	r := chi.NewRouter()
	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		io.WriteString(rw, "Hello, world!")
	})
	// r.Get("/user", HandleGetUser)
	r.Post("/user/new", HandleCreateUser)
	r.Delete("/user", HandleDeleteUser)

	n.UseHandler(r)

	log.Info("Listening on :3001...")
	http.ListenAndServe(":3001", n)
}
