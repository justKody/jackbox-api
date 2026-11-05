package v1

import (
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRoutes(router *mux.Router) {
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods(http.MethodGet)
}
