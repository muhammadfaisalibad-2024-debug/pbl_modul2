package model

type Prestasi struct {
	ID        int    `json:"id"`
	StudentID int    `json:"student_id"`
	Nama      string `json:"nama"`
	Juara     int    `json:"juara"`
}
