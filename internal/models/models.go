package models



type Level struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Order       int        `json:"order"`
	Skills      string     `json:"skills"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	Exercises   []Exercise `json:"exercises,omitempty"`
}

type Exercise struct {
	ID              int    `json:"id"`
	LevelID         int    `json:"level_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	SkillCategory   string `json:"skill_category"`
	Difficulty      string `json:"difficulty"`
	Equipment       string `json:"equipment"`
	DurationMinutes int    `json:"duration_minutes"`
	YoutubeURL      string `json:"youtube_url"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type ExerciseCreate struct {
	LevelID         int    `json:"level_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	SkillCategory   string `json:"skill_category"`
	Difficulty      string `json:"difficulty"`
	Equipment       string `json:"equipment"`
	DurationMinutes int    `json:"duration_minutes"`
	YoutubeURL      string `json:"youtube_url"`
}
