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

	if err := seedExercises(); err != nil {
		return err
	}
	return nil
}

func seedExercises() error {
	exercises := []models.ExerciseCreate{
		// Level 1 — Water Competency
		{LevelID: 1, Name: "Freestyle Technique Drill", Description: "50m freestyle focusing on body rotation, bilateral breathing, and high elbow catch. Coach observes from deck giving real-time feedback.", SkillCategory: "swimming", Difficulty: "beginner", Equipment: "kickboard, pull buoy", DurationMinutes: 15, YoutubeURL: "https://www.youtube.com/watch?v=s20VLN7ACWI"},
		{LevelID: 1, Name: "Eggbeater Stability Challenge", Description: "Players maintain eggbeater kick while holding ball overhead with both hands for 30s intervals. Progress to single hand, then catching passes while stationary.", SkillCategory: "conditioning", Difficulty: "beginner", Equipment: "water polo ball", DurationMinutes: 10, YoutubeURL: "https://www.youtube.com/watch?v=q3xC0nF65Dc"},
		{LevelID: 1, Name: "Ball Familiarization Circuit", Description: "5 stations: grip squeeze, toss-and-catch, two-hand overhead throw, dribble 10m, scoop pickup from water. 2 min per station, rotate.", SkillCategory: "ball_handling", Difficulty: "beginner", Equipment: "water polo balls (1 per player), cones", DurationMinutes: 12, YoutubeURL: ""},
		// Level 2 — Passing & Movement
		{LevelID: 2, Name: "Partner Passing Progression", Description: "Pairs 3m apart: 10 static passes each hand, then move to 5m, then passing while treading. Focus on wrist snap, ball spin, and finger placement on release.", SkillCategory: "passing", Difficulty: "beginner", Equipment: "water polo balls", DurationMinutes: 15, YoutubeURL: "https://www.youtube.com/watch?v=FpJEnMbR9xY"},
		{LevelID: 2, Name: "Moving Catch & Release", Description: "Player swims head-up, receives pass from coach, immediately passes to target without stopping. Alternate sides. Emphasize one fluid motion.", SkillCategory: "passing", Difficulty: "intermediate", Equipment: "water polo balls, floating targets", DurationMinutes: 12, YoutubeURL: ""},
		{LevelID: 2, Name: "Dribble Relay Race", Description: "Teams of 3. Dribble head-up across pool, tag teammate. Ball must stay within arm's reach. Emphasize speed with control.", SkillCategory: "swimming", Difficulty: "beginner", Equipment: "water polo balls, lane ropes", DurationMinutes: 10, YoutubeURL: ""},
		// Level 3 — Pressure & Tactics
		{LevelID: 3, Name: "Faking Drill (Short & Medium)", Description: "Player at 5m: perform short fake (ball doesn't leave hand), then medium fake (arm goes back). Defender reacts. 10 reps each type, alternate attacker/defender.", SkillCategory: "shooting", Difficulty: "intermediate", Equipment: "water polo balls, goal", DurationMinutes: 15, YoutubeURL: "https://www.youtube.com/watch?v=KcCN8F3YNAM"},
		{LevelID: 3, Name: "3v2 Man-Up Basic", Description: "Offensive triangle vs 2 defenders. Attack must complete 3 passes before shooting. Defenders communicate and rotate. Teach basic man-up positioning.", SkillCategory: "tactics", Difficulty: "intermediate", Equipment: "water polo balls, goal, caps", DurationMinutes: 20, YoutubeURL: ""},
		{LevelID: 3, Name: "Ball Shielding Under Pressure", Description: "Attacker holds ball, defender applies light body pressure from behind. Attacker maintains possession using body positioning for 15s. Rotate roles.", SkillCategory: "ball_handling", Difficulty: "intermediate", Equipment: "water polo balls", DurationMinutes: 10, YoutubeURL: ""},
		// Level 4 — Individual Dominance
		{LevelID: 4, Name: "1v1 Drive to Outside Lane", Description: "Attacker at 5m, defender on hip. Attacker must win outside position and receive ball for shot. Focus on swim speed, body contact, and positioning.", SkillCategory: "swimming", Difficulty: "advanced", Equipment: "water polo balls, goal", DurationMinutes: 15, YoutubeURL: ""},
		{LevelID: 4, Name: "Press Defense Fundamentals", Description: "Defender practices pressing technique: hand on ball, hips low, feet moving. Attacker tries to turn and shoot. 30s rounds, switch roles.", SkillCategory: "defense", Difficulty: "advanced", Equipment: "water polo balls, goal", DurationMinutes: 15, YoutubeURL: "https://www.youtube.com/watch?v=VrX48GP3bOU"},
		{LevelID: 4, Name: "Contact Shooting", Description: "Shooter receives ball at 4m with defender on back. Must create space and shoot within 3 seconds. Vary defender pressure from 50% to 100%.", SkillCategory: "shooting", Difficulty: "advanced", Equipment: "water polo balls, goal, caps", DurationMinutes: 15, YoutubeURL: ""},
		// Level 5 — Collective Play
		{LevelID: 5, Name: "Pick and Roll Execution", Description: "Center sets pick on wing defender, wing drives to 3m. Practice timing: too early = foul, too late = no advantage. Run from both sides.", SkillCategory: "tactics", Difficulty: "advanced", Equipment: "water polo balls, goal, caps", DurationMinutes: 20, YoutubeURL: ""},
		{LevelID: 5, Name: "6v5 Man-Up Rotation Drill", Description: "Full man-up set: 4-2 positioning. Practice ball movement around the arc. Dry pass to post when defense collapses. 5 possessions, then rotate players.", SkillCategory: "tactics", Difficulty: "advanced", Equipment: "water polo balls, goal, caps", DurationMinutes: 25, YoutubeURL: "https://www.youtube.com/watch?v=k6CZm8b1wWo"},
		{LevelID: 5, Name: "Defensive Switching Communication", Description: "3v3 half-court. Offense picks and rolls, defense must call switches loudly. Coach stops play if no verbal communication. Build habit of calling 'switch' and 'stay'.", SkillCategory: "defense", Difficulty: "advanced", Equipment: "water polo balls, goal, caps", DurationMinutes: 15, YoutubeURL: ""},
		// Level 6 — High-Level Skills
		{LevelID: 6, Name: "Advanced Faking Combo", Description: "Long fake → drive → backhand shot. Then: double fake → step-out → power shot. 5 reps each combo against active defender at 75% pressure.", SkillCategory: "shooting", Difficulty: "elite", Equipment: "water polo balls, goal", DurationMinutes: 20, YoutubeURL: ""},
		{LevelID: 6, Name: "Center Play: Seal & Roll", Description: "Center receives ball at 2m with back to goal. Practice seal (pin defender), roll to shooting position, backhand or sweep shot. Alternate left/right side.", SkillCategory: "tactics", Difficulty: "elite", Equipment: "water polo balls, goal, caps", DurationMinutes: 20, YoutubeURL: ""},
		// Level 7 — Competition Mastery
		{LevelID: 7, Name: "Full 6v6 Tactical Scrimmage", Description: "Game-speed scrimmage with specific constraints: must execute 2 picks per possession, all goals must come from center or drive. Coach stops play for teaching moments.", SkillCategory: "tactics", Difficulty: "elite", Equipment: "water polo balls, goals, caps, shot clock", DurationMinutes: 30, YoutubeURL: ""},
		{LevelID: 7, Name: "Fatigue Decision-Making Drill", Description: "Players sprint 50m, then immediately enter 3v2 situation. Must make correct pass/shoot decision under fatigue within 5 seconds. Simulates late-game scenarios.", SkillCategory: "conditioning", Difficulty: "elite", Equipment: "water polo balls, goal, caps", DurationMinutes: 20, YoutubeURL: ""},
	}

	for _, ex := range exercises {
		_, err := DB.Exec(
			`INSERT INTO exercises (level_id, name, description, skill_category, difficulty, equipment, duration_minutes, youtube_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			ex.LevelID, ex.Name, ex.Description, ex.SkillCategory, ex.Difficulty, ex.Equipment, ex.DurationMinutes, ex.YoutubeURL,
		)
		if err != nil {
			return fmt.Errorf("seed exercise %s: %w", ex.Name, err)
		}
	}
	log.Printf("Seeded %d exercises across all levels", len(exercises))
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
