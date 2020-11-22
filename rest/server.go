package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/paulbellamy/snippets/snippets"
)

type Server struct {
	http.Handler
	baseUrl  string
	snippets snippets.Store
}

func NewServer(stdout io.Writer, baseUrl string, snippetStore snippets.Store) *Server {
	r := mux.NewRouter()
	s := &Server{
		Handler:  handlers.LoggingHandler(stdout, r),
		baseUrl:  baseUrl,
		snippets: snippetStore,
	}

	r.HandleFunc("/snippets", s.postSnippet).Methods("POST")
	r.HandleFunc("/snippets/{name}", s.getSnippet).Methods("GET")

	return s
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

	stored, err := s.snippets.Store(input.Name, input.Snippet, input.ExpiresIn*time.Second)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := s.writeSnippet(w, stored); err != nil {
		log.Println("Error:", err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getSnippet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	value, err := s.snippets.Load(vars["name"])
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	if value == nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}

	if err := s.writeSnippet(w, value); err != nil {
		log.Println("Error:", err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

type snippetResponse struct {
	ExpiresAt time.Time `json:"expires_at"`
	Name      string    `json:"name"`
	Snippet   string    `json:"snippet"`
	Url       string    `json:"url"`
}

func (s *Server) writeSnippet(w io.Writer, snip *snippets.Snippet) error {
	return json.NewEncoder(w).Encode(snippetResponse{
		ExpiresAt: snip.ExpiresAt,
		Name:      snip.Name,
		Snippet:   snip.Snippet,
		Url:       s.baseUrl + "/snippets/" + snip.Name,
	})
}
