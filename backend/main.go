package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"transit-ni/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET"},
	}))

	r.Get("/api/stops", handlers.GetStops)
	r.Get("/api/stops/{stopID}", handlers.GetStop)
	r.Get("/api/routes", handlers.GetRoutes)
	r.Get("/api/plan", handlers.PlanTrip)

	log.Println("backend running on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
