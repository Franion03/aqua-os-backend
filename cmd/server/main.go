package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/franion03/aqua-os-backend/internal/auth"
	"github.com/franion03/aqua-os-backend/internal/db"
	"github.com/franion03/aqua-os-backend/internal/handlers"
)

func main() {
	// ── Database ──────────────────────────────────────────────────
	if err := db.Init("aquaos.db"); err != nil {
		log.Fatalf("Database init failed: %v", err)
	}
	defer db.Close()

	// ── Router ────────────────────────────────────────────────────
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// ── Routes ────────────────────────────────────────────────────
	r.Get("/api/health", handlers.Health)

	// Auth
	r.Post("/api/auth/login", handlers.Login)

	// Levels (public read-only)
	r.Route("/api/levels", func(r chi.Router) {
		r.Get("/", handlers.ListLevels)
		r.Get("/{id}", handlers.GetLevel)
	})

	// Exercises — read public, write protected
	r.Route("/api/exercises", func(r chi.Router) {
		r.Get("/", handlers.ListExercises)

		// Protected: require Admin JWT
		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware)
			r.Use(auth.RequireRole(auth.AdminRole))
			r.Post("/", handlers.AddExercise)
			r.Delete("/{id}", handlers.DeleteExercise)
		})
	})

	// ── Start ────────────────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("AquaOS Backend starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
