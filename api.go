package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr  string
	store Store
}

func NewAPIServer(addr string, store Store) *APIServer {
	return &APIServer{
		addr:  addr,
		store: store,
	}
}

func (s *APIServer) Serve() {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	us := NewUSerService(s.store)
	us.RegisterRoutes(subRouter)

	ts := NewTaskService(s.store)
	ts.RegisterRoutes(subRouter)

	log.Println("Starting server at ", s.addr)
	log.Fatal(http.ListenAndServe(s.addr, subRouter))
}
