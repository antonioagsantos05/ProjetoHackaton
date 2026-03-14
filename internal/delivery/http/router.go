package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/fiap-x/video-processor/internal/delivery/http/middleware"
	// "github.com/fiap-x/video-processor/internal/delivery/http/handlers"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.AuthMiddleware) // Adiciona JWT global ou em sub-routers
	r.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "UP"}`))
	}).Methods("GET")

	return r
}
