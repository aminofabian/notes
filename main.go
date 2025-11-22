package main

import (
	"net/http"

	"github.com/aminofabian/notes/controllers"
	"github.com/aminofabian/notes/middleware"
	"github.com/gorilla/mux"
)

func main() {

	// Initialize router
	r := mux.NewRouter()

	// Apply CORS middleware to all routes
	r.Use(middleware.EnableCORS)

	// Routes
	r.HandleFunc("/",
		controllers.Hello,
	)

	r.HandleFunc("/notes",
		controllers.GetNotes,
	).Methods("POST", "OPTIONS")

	// Start server
	http.ListenAndServe(":8080", r)

}
