package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"

	"github.com/franion03/aqua-os-backend/internal/auth"
	"github.com/franion03/aqua-os-backend/internal/db"
	"github.com/franion03/aqua-os-backend/internal/handlers"
)

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/health", handlers.Health)
	r.Get("/api/levels", handlers.ListLevels)
	r.Route("/api/exercises", func(r chi.Router) {
		r.Get("/", handlers.ListExercises)
		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware)
			r.Use(auth.RequireRole(auth.AdminRole))
			r.Post("/", handlers.AddExercise)
		})
	})
	return r
}

// generateTestToken creates a token that passes the auth middleware validation
// (includes iss and aud claims that the middleware checks).
func generateTestToken() string {
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		jwtKey = "CHANGE_ME_TO_A_LONG_RANDOM_SECRET_KEY_AT_LEAST_32_CHARS"
	}
	jwtIssuer := os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		jwtIssuer = "AquaOsCalendar"
	}
	jwtAud := os.Getenv("JWT_AUDIENCE")
	if jwtAud == "" {
		jwtAud = "AquaOsCalendar"
	}

	claims := jwt.MapClaims{
		"role": auth.AdminRole,
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  jwtIssuer,
		"aud":  jwtAud,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(jwtKey))
	return s
}

func TestMain(m *testing.M) {
	// Open in-memory DB with timestamp parsing enabled
	conn, err := sql.Open("sqlite3", "file::memory:?_foreign_keys=on&_loc=auto")
	if err != nil {
		panic(err)
	}
	db.DB = conn

	// Create tables and seed
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS levels (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			"order" INTEGER NOT NULL,
			skills TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT (datetime('now')),
			updated_at TIMESTAMP DEFAULT (datetime('now'))
		);
		CREATE TABLE IF NOT EXISTS exercises (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			level_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			skill_category TEXT NOT NULL DEFAULT 'general',
			difficulty TEXT NOT NULL DEFAULT 'beginner',
			equipment TEXT DEFAULT '',
			duration_minutes INTEGER DEFAULT 15,
			youtube_url TEXT DEFAULT '',
			created_at TEXT DEFAULT (datetime('now')),
			updated_at TEXT DEFAULT (datetime('now')),
			FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE CASCADE
		);
		INSERT INTO levels (id, name, description, "order", skills) VALUES
		(1, 'Water Competency', 'Basic water confidence', 1, 'swimming||eggbeater');
	`)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	conn.Close()
	os.Exit(code)
}

func TestHealthEndpoint(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestListLevels(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/levels", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]json.RawMessage
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := resp["levels"]; !ok {
		t.Fatal("response missing 'levels' key")
	}
	// Verify it's an array
	var levels []interface{}
	if err := json.Unmarshal(resp["levels"], &levels); err != nil {
		t.Fatalf("levels is not an array: %v", err)
	}
}

func TestAddExercise(t *testing.T) {
	r := setupRouter()
	body, _ := json.Marshal(map[string]interface{}{
		"level_id":    1,
		"name":        "Eggbeater Drill",
		"description": "Basic eggbeater practice",
	})

	token := generateTestToken()

	req := httptest.NewRequest(http.MethodPost, "/api/exercises", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestListExercises(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/exercises", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]json.RawMessage
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := resp["exercises"]; !ok {
		t.Fatal("response missing 'exercises' key")
	}
	var exercises []interface{}
	if err := json.Unmarshal(resp["exercises"], &exercises); err != nil {
		t.Fatalf("exercises is not an array: %v", err)
	}
}
