package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	http.Handler
	baseUrl string
}

func NewServer(baseUrl string) *Server {
	r := mux.NewRouter()
	s := &Server{
		Handler: r,
		baseUrl: baseUrl,
	}

	r.HandleFunc("/snippets", s.postSnippet).Methods("POST")

	return s
}

type snippetRequest struct {
	ExpiresIn time.Duration `json:"expires_in"`
	Name      string        `json:"name"`
	Snippet   string        `json:"snippet"`
}

type snippetResponse struct {
	ExpiresAt time.Time `json:"expires_at"`
	Name      string    `json:"name"`
	Snippet   string    `json:"snippet"`
	Url       string    `json:"url"`
}

func (s *Server) postSnippet(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// TODO: Limit body size.
	var request snippetRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: validate the input
	// - name should not be blank
	// - expiresAt should be postive (and less than some max)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(snippetResponse{
		ExpiresAt: time.Now().Add(request.ExpiresIn * time.Second),
		Name:      request.Name,
		Snippet:   request.Snippet,
		Url:       s.baseUrl + "/snippets/" + request.Name,
	})
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}
