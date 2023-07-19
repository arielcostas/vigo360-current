package internal

import (
	"fmt"
	"net/http"
	"time"
	"vigo360.es/new/internal/templates"

	"vigo360.es/new/internal/messages"
)

type errorResponse struct {
	Rid          string
	ErrorCode    int
	Message      string
	RequestedUrl string
	Time         string
}

func (s *Server) handleError(r *http.Request, w http.ResponseWriter, status int, message messages.ErrorMessage) {
	w.WriteHeader(status)
	rid := r.Context().Value(ridContextKey("rid")).(string)

	err := templates.Render(w, "_error.html", &errorResponse{
		Rid:          rid,
		ErrorCode:    status,
		Message:      fmt.Sprintf("%s", message),
		RequestedUrl: r.URL.String(),
		Time:         time.Now().Format("2006-01-02 15:04:05"),
	})

	if err != nil {
		_ = templates.Render(w, "_error.html", &errorResponse{
			Rid:          rid,
			ErrorCode:    status,
			Message:      fmt.Sprintf("%s", message),
			RequestedUrl: r.URL.String(),
			Time:         time.Now().Format("2006-01-02 15:04:05"),
		})
		fmt.Println(rid, err)
	}
}

func (s *Server) handleJsonError(r *http.Request, w http.ResponseWriter, status int, message messages.ErrorMessage) {
	w.WriteHeader(status)
	rid := r.Context().Value(ridContextKey("rid")).(string)
	_, err := fmt.Fprintf(w, `{ "error": "%s", "rid": "%s" }`, message, rid)
	if err != nil {
		_, _ = fmt.Fprintf(w, `{ "error": "%s", "rid": "%s" }`, messages.ErrorFatal, rid)
		fmt.Println(rid, err)
	}
}
