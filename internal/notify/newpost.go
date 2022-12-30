package notify

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func NotifyNewPost(id string, titulo string, autor string) error {
	// Create an HTTP client and send a request to the URL defined by the "POST_NOTIFICATION_URL" environment variable
	var url = os.Getenv("POST_NOTIFICATION_URL")
	var content = fmt.Sprintf(`{"id": "%s", "titulo": "%s", "autor": "%s"}`, id, titulo, autor)

	_, err := http.DefaultClient.Post(
		os.Getenv(url),
		"application/json",
		strings.NewReader(
			content,
		),
	)

	return err
}
