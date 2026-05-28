package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/franion03/aqua-os-backend/internal/auth"
	"github.com/franion03/aqua-os-backend/internal/db"
	"github.com/franion03/aqua-os-backend/internal/models"
)

// ── Health ──────────────────────────────────────────────────────────

func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "aqua-os-backend",
	})
}

// ── Levels ──────────────────────────────────────────────────────────

func ListLevels(w http.ResponseWriter, r *http.Request) {
	levels, err := db.GetAllLevels()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if levels == nil {
		levels = []models.Level{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"levels": levels})
}

func GetLevel(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid level id")
		return
	}

	level, err := db.GetLevel(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if level == nil {
		writeError(w, http.StatusNotFound, "level not found")
		return
	}
	if level.Exercises == nil {
		level.Exercises = []models.Exercise{}
	}
	writeJSON(w, http.StatusOK, level)
}

// ── Exercises ───────────────────────────────────────────────────────

func ListExercises(w http.ResponseWriter, r *http.Request) {
	var levelID *int
	if val := r.URL.Query().Get("level_id"); val != "" {
		id, err := strconv.Atoi(val)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid level_id")
			return
		}
		levelID = &id
	}

	exercises, err := db.GetExercises(levelID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if exercises == nil {
		exercises = []models.Exercise{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"exercises": exercises})
}

func AddExercise(w http.ResponseWriter, r *http.Request) {
	var req models.ExerciseCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Name == "" || req.Description == "" || req.LevelID == 0 {
		writeError(w, http.StatusBadRequest, "name, description, and level_id are required")
		return
	}

	exercise, err := db.AddExercise(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, exercise)
}

func DeleteExercise(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid exercise id")
		return
	}
	if err := db.DeleteExercise(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ── Auth ────────────────────────────────────────────────────────────

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminUser == "" {
		adminUser = "CHANGE_ME"
	}
	if adminPass == "" {
		adminPass = "CHANGE_ME"
	}

	if req.Username != adminUser || req.Password != adminPass {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := auth.GenerateToken(auth.AdminRole)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token})
}

// ── Helpers ─────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
