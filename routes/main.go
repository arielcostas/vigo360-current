package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func TestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:06"), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func InitRouter() *mux.Router {
	InitDB()

	router := mux.NewRouter().StrictSlash(true)
	router.Use(TestMiddleware)
	router.HandleFunc("/", HomePage)
	router.HandleFunc("/admin/login", AdminLogin).Methods("GET", "POST")
	return router
}

func ValidPassword(password string, hash string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if res != nil {
		println(res.Error())
		return false
	}
	return true
}
