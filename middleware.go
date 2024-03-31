package main

import (
	"net/http"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/urfave/negroni"
)

func LogMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lrw := negroni.NewResponseWriter(rw)

	next(lrw, r)

	logRequest(r.URL.Path, r.Method, lrw.Status())
}

func logRequest(path string, method string, status int) {
	styles := log.DefaultStyles()
	color := 2
	if status >= 400 && status < 500 {
		color = 4
	}
	if status >= 500 {
		color = 5
	}
	if status >= 300 && status < 400 {
		color = 3
	}

	styles.Values["status"] = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(color))
	logger := log.New(os.Stdout)
	logger.SetStyles(styles)
	logger.Info("request", "method", method, "path", path, "status", status)
}
