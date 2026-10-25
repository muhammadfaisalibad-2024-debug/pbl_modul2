package model

import "time"

// Portfolio adalah representasi portofolio mahasiswa yang dikenal seluruh aplikasi.
// Sesuai aturan arsitektur, package model tidak mengimpor driver MongoDB.
// Seluruh ID direpresentasikan sebagai string.
type Portfolio struct {
	ID           string        `json:"id"`
	StudentID    int           `json:"student_id"`
	Headline     string        `json:"headline"`
	Bio          string        `json:"bio"`
	Projects     []Project     `json:"projects"`
	Certificates []Certificate `json:"certificates"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type Project struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Technologies []string `json:"technologies"`
	URL          string   `json:"url,omitempty"`
	Year         int      `json:"year"`
}

type Certificate struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Issuer    string    `json:"issuer"`
	Date      time.Time `json:"date"`
}

type PutPortfolioRequest struct {
	Headline     string             `json:"headline" validate:"required,min=3,max=100"`
	Bio          string             `json:"bio" validate:"max=1000"`
	Projects     []ProjectInput     `json:"projects" validate:"max=20,dive"`
	Certificates []CertificateInput `json:"certificates" validate:"max=50,dive"`
}

type ProjectInput struct {
	Title        string   `json:"title" validate:"required,min=2,max=100"`
	Technologies []string `json:"technologies" validate:"max=15,dive,min=1,max=30"`
	URL          string   `json:"url,omitempty" validate:"omitempty,url"`
	Year         int      `json:"year" validate:"min=2000,max=2100"`
}

type CertificateInput struct {
	Name   string `json:"name" validate:"required,min=2,max=100"`
	Issuer string `json:"issuer" validate:"required,min=2,max=100"`
	Date   string `json:"date" validate:"required"`
}

type AddProjectRequest = ProjectInput

type TechCount struct {
	Technology string `json:"technology"`
	Total      int    `json:"total"`
}
