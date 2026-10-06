package model

import "time"

type Course struct {
	ID        int       `json:"id"`
	KodeMK    string    `json:"kode_mk" validate:"required,max=20"`
	NamaMK    string    `json:"nama_mk" validate:"required,max=120"`
	SKS       int       `json:"sks" validate:"required,min=1"`
	Semester  int       `json:"semester" validate:"required,min=1"`
	Kuota     int       `json:"kuota" validate:"required,gte=0"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type CourseListItem struct {
	ID        int    `json:"id"`
	KodeMK    string `json:"kode_mk"`
	NamaMK    string `json:"nama_mk"`
	SKS       int    `json:"sks"`
	Semester  int    `json:"semester"`
	Kuota     int    `json:"kuota"`
	Terisi    int    `json:"terisi"`
	SisaKuota int    `json:"sisa_kuota"`
}

type CourseQuery struct {
	Semester  int
	Search    string
	Available bool
}
