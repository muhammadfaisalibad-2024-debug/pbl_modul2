package model

import "time"

type Student struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id,omitempty"`
	NIM         string     `json:"nim" validate:"required,len=12,numeric"`
	Nama        string     `json:"nama" validate:"required,max=120"`
	Prodi       string     `json:"prodi" validate:"required,max=100"`
	Angkatan    int        `json:"angkatan" validate:"required,min=1900"`
	IPKTerakhir float64    `json:"ipk_terakhir" validate:"gte=0,lte=4"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
}

type StudentListItem struct {
	ID          int     `json:"id"`
	NIM         string  `json:"nim"`
	Nama        string  `json:"nama"`
	Prodi       string  `json:"prodi"`
	Angkatan    int     `json:"angkatan"`
	IPKTerakhir float64 `json:"ipk_terakhir"`
}

type StudentDetailResponse struct {
	ID          int                  `json:"id"`
	NIM         string               `json:"nim"`
	Nama        string               `json:"nama"`
	Prodi       string               `json:"prodi"`
	Angkatan    int                  `json:"angkatan"`
	IPKTerakhir float64              `json:"ipk_terakhir"`
	TotalSKS    int                  `json:"total_sks"`
	BatasSKS    int                  `json:"batas_sks"`
	MataKuliah  []EnrolledCourseItem `json:"mata_kuliah"`
}

type CreateStudentRequest struct {
	NIM         string   `json:"nim" validate:"required,len=12,numeric"`
	Nama        string   `json:"nama" validate:"required,min=2,max=120"`
	Email       string   `json:"email" validate:"required,email,max=120"`
	Prodi       string   `json:"prodi" validate:"required,min=2,max=100"`
	Angkatan    int      `json:"angkatan" validate:"required,min=1900"`
	IPKTerakhir *float64 `json:"ipk_terakhir,omitempty" validate:"omitempty,gte=0,lte=4"`
}

type UpdateStudentRequest struct {
	Nama        string   `json:"nama" validate:"required,min=2,max=120"`
	Prodi       string   `json:"prodi" validate:"required,min=2,max=100"`
	Angkatan    int      `json:"angkatan" validate:"required,min=1900"`
	IPKTerakhir *float64 `json:"ipk_terakhir,omitempty" validate:"omitempty,gte=0,lte=4"`
}

type StudentListQuery struct {
	Page     int
	PerPage  int
	Prodi    string
	Angkatan int
	Search   string
	Sort     string
}

func CalculateBatasSKS(ipk float64) int {
	if ipk >= 3.00 {
		return 24
	} else if ipk >= 2.50 {
		return 21
	}
	return 18
}
