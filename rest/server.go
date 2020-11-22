package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/paulbellamy/snippets/snippets"
)

type Server struct {
	http.Handler
	baseUrl string
	s       snippets.Store
}

func NewServer(baseUrl string, snippetStore snippets.Store) *Server {
	r := mux.NewRouter()
	s := &Server{
		Handler: r,
		baseUrl: baseUrl,
		s:       snippetStore,
	}

	r.HandleFunc("/snippets", s.postSnippet).Methods("POST")

	return s
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
	var input snippets.Snippet
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: validate the input
	// - name should not be blank
	// - expiresIn should be postive (and less than some max)

	stored, err := s.s.Store(input.Name, input.Snippet, input.ExpiresIn)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(snippetResponse{
		ExpiresAt: stored.ExpiresAt,
		Name:      stored.Name,
		Snippet:   stored.Snippet,
		Url:       s.baseUrl + "/snippets/" + stored.Name,
	})
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}
