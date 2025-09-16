package server

import (
	"fmt"
	"github.com/supercute/spx-graph/internal/graph"
	"net/http"
)

type Server struct {
	generator *graph.Generator
	port      int
}

func New(generator *graph.Generator, port int) *Server {
	return &Server{
		generator: generator,
		port:      port,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.handleGraph)
	http.HandleFunc("/favicon.ico", s.handleFavicon)

	addr := fmt.Sprintf(":%d", s.port)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleGraph(w http.ResponseWriter, r *http.Request) {
	// Generate SVG Graph
	svg, err := s.generator.GenerateSVG()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate graph: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate HTML with SVG
	html := s.generator.GenerateHTML(svg)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func (s *Server) handleFavicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
