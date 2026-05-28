package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/franion03/aqua-os-backend/internal/models"
)

var DB *sql.DB

var levelSeeds = []models.Level{
	{ID: 1, Name: "Water Competency & Basic Ball Handling", Order: 1,
		Description: "Establish fundamental water confidence and comfort. Players master proper freestyle technique, basic eggbeater kick, vertical positioning, and basic ball handling.",
		Skills:      "Proper freestyle swimming technique with correct body alignment and breathing||Strong basic eggbeater kick for stationary stability||Vertical positioning and controlled vertical jumps||Basic ball handling (catching, gripping, and two-handed throwing)||Head-up swimming and simple dribbling||Safe water entries, gliding, and streamline position"},
	{ID: 2, Name: "Controlled Passing & Continuous Movement", Order: 2,
		Description: "Develop consistent passing and catching technique while learning to move effectively with the ball.",
		Skills:      "Proper passing technique (wrist snap, ball spin, and accurate release)||Reliable catching with one and two hands while stationary and moving||Core stability and balance during ball handling||Continuous movement while dribbling with head up||Moving catch-and-release in one fluid motion||Basic eggbeater while performing arm actions (passing/sculling)"},
	{ID: 3, Name: "Handling Pressure & Basic Tactics", Order: 3,
		Description: "Introduce light defensive pressure while refining passing and shooting accuracy. Basic tactical concepts introduced.",
		Skills:      "Passing and catching under light pressure||Basic faking (short and medium fakes)||Ball shielding and protection with the body||Proper shooting technique (wrist shot and push shot)||Simple tactical concepts: man-up positioning and basic zone defense||Swim-through movements to create space||Maintaining verticality and stability when contested"},
	{ID: 4, Name: "Individual Dominance & Defensive Fundamentals", Order: 4,
		Description: "Build individual 1v1 skills on both offense and defense.",
		Skills:      "Winning the outside lane in 1v1 situations||Strong defensive positioning, footwork, and press defense||Contact shooting (shooting while defended)||Effective ball protection in physical duels||Basic transition play (offense to defense and vice versa)||Proper defensive posture (hand up, hips low, controlled distance)||Advanced eggbeater for explosive vertical jumps under pressure"},
	{ID: 5, Name: "Collective Play & Tactical Versatility", Order: 5,
		Description: "Develop team structure and coordinated play. Positional roles introduced.",
		Skills:      "Pick and roll execution (timing and spacing)||Advanced man-up patterns and rotations||Defensive player switching and communication||Positional role fundamentals (Center, Center-Back, Wing, Goalie)||M-Zone defensive structure and rotation timing||Tactical patience and structured team attacks||Understanding of basic counterattack timing"},
	{ID: 6, Name: "High-Level Skills & Adaptive Teamplay", Order: 6,
		Description: "Refine advanced technical skills and develop game intelligence.",
		Skills:      "Advanced faking (long fakes, shoulder fakes, double fakes, fake-and-drive)||Advanced defensive techniques (steals, shot blocking, fronting the center)||Center position mastery (seal & roll, backhand shots)||Reading opponent defenses and adapting in real time||Complex transition play (counter, reset, or press decisions)||High-level ball control under heavy pressure||Fluid combination of fakes, drives, and shots"},
	{ID: 7, Name: "High Performance & Competition Mastery", Order: 7,
		Description: "Prepare players for competitive performance at the highest level.",
		Skills:      "Tactical leadership and play-calling on the field||Full game system execution (complex man-up, double picks, drive rotations)||Elite decision-making under fatigue and high pressure||Reading complex game situations and opponent tendencies||Maintaining technical precision and speed at competition intensity||Advanced pressing variations and defensive adjustments||Complete transition mastery in all game phases"},
}

func Init(dbPath string) error {
	if dbPath == "" {
		dbPath = filepath.Join(filepath.Dir(os.Args[0]), "aquaos.db")
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	if err = createTables(); err != nil {
		return err
	}

	if err = seedLevels(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS levels (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			"order" INTEGER NOT NULL,
			skills TEXT NOT NULL,
			created_at TEXT DEFAULT (datetime('now')),
			updated_at TEXT DEFAULT (datetime('now'))
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
	`)
	return err
}

func seedLevels() error {
	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM levels").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	for _, level := range levelSeeds {
		_, err := DB.Exec(
			`INSERT INTO levels (id, name, description, "order", skills) VALUES (?, ?, ?, ?, ?)`,
			level.ID, level.Name, level.Description, level.Order, level.Skills,
		)
		if err != nil {
			return fmt.Errorf("seed level %d: %w", level.ID, err)
		}
	}
	log.Println("Seeded 7 levels")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

// ── Level queries ───────────────────────────────────────────────────

func GetAllLevels() ([]models.Level, error) {
	rows, err := DB.Query(`SELECT id, name, description, "order", skills, created_at, updated_at FROM levels ORDER BY "order"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []models.Level
	for rows.Next() {
		var l models.Level
		if err := rows.Scan(&l.ID, &l.Name, &l.Description, &l.Order, &l.Skills, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		levels = append(levels, l)
	}
	return levels, nil
}

func GetLevel(id int) (*models.Level, error) {
	var l models.Level
	err := DB.QueryRow(
		`SELECT id, name, description, "order", skills, created_at, updated_at FROM levels WHERE id = ?`, id,
	).Scan(&l.ID, &l.Name, &l.Description, &l.Order, &l.Skills, &l.CreatedAt, &l.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`SELECT id, level_id, name, description, skill_category, difficulty, equipment, duration_minutes, youtube_url, created_at, updated_at FROM exercises WHERE level_id = ? ORDER BY id DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Exercise
		if err := rows.Scan(&e.ID, &e.LevelID, &e.Name, &e.Description, &e.SkillCategory, &e.Difficulty, &e.Equipment, &e.DurationMinutes, &e.YoutubeURL, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		l.Exercises = append(l.Exercises, e)
	}
	return &l, nil
}

// ── Exercise queries ────────────────────────────────────────────────

func AddExercise(req models.ExerciseCreate) (*models.Exercise, error) {
	if req.SkillCategory == "" {
		req.SkillCategory = "general"
	}
	if req.Difficulty == "" {
		req.Difficulty = "beginner"
	}
	if req.DurationMinutes == 0 {
		req.DurationMinutes = 15
	}

	result, err := DB.Exec(
		`INSERT INTO exercises (level_id, name, description, skill_category, difficulty, equipment, duration_minutes, youtube_url)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		req.LevelID, req.Name, req.Description, req.SkillCategory, req.Difficulty, req.Equipment, req.DurationMinutes, req.YoutubeURL,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	var e models.Exercise
	err = DB.QueryRow(
		`SELECT id, level_id, name, description, skill_category, difficulty, equipment, duration_minutes, youtube_url, created_at, updated_at FROM exercises WHERE id = ?`, id,
	).Scan(&e.ID, &e.LevelID, &e.Name, &e.Description, &e.SkillCategory, &e.Difficulty, &e.Equipment, &e.DurationMinutes, &e.YoutubeURL, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func GetExercises(levelID *int) ([]models.Exercise, error) {
	var rows *sql.Rows
	var err error

	if levelID != nil {
		rows, err = DB.Query(
			`SELECT e.id, e.level_id, e.name, e.description, e.skill_category, e.difficulty, e.equipment, e.duration_minutes, e.youtube_url, e.created_at, e.updated_at
			 FROM exercises e WHERE e.level_id = ? ORDER BY e.id DESC`, *levelID,
		)
	} else {
		rows, err = DB.Query(
			`SELECT e.id, e.level_id, e.name, e.description, e.skill_category, e.difficulty, e.equipment, e.duration_minutes, e.youtube_url, e.created_at, e.updated_at
			 FROM exercises e ORDER BY e.id DESC`,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var e models.Exercise
		if err := rows.Scan(&e.ID, &e.LevelID, &e.Name, &e.Description, &e.SkillCategory, &e.Difficulty, &e.Equipment, &e.DurationMinutes, &e.YoutubeURL, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}
	return exercises, nil
}

func DeleteExercise(id int) error {
	_, err := DB.Exec("DELETE FROM exercises WHERE id = ?", id)
	return err
}
