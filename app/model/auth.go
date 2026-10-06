package model

import "time"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,username"`
	Email    string `json:"email" validate:"required,email,max=120"`
	Password string `json:"password" validate:"required,min=8,max=72,strongpassword"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	User         *UserData `json:"user,omitempty"`
}

type UserData struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type StudentProfile struct {
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	Prodi    string `json:"prodi"`
	Angkatan int    `json:"angkatan"`
}

type MeResponse struct {
	ID      int             `json:"id"`
	Email   string          `json:"email"`
	Role    string          `json:"role"`
	Student *StudentProfile `json:"student,omitempty"`
}

type RefreshToken struct {
	ID        int64
	UserID    int
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type AuthUser struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
