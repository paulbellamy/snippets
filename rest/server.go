package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	http.Handler
}

func NewServer() *Server {
	r := mux.NewRouter()
	s := &Server{
		Handler: r,
	}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	return s
}
