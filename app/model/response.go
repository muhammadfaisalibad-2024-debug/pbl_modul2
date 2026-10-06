package model

type WebResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Success   bool                `json:"success"`
	Message   string              `json:"message"`
	Errors    map[string][]string `json:"errors,omitempty"`
	RequestID string              `json:"request_id,omitempty"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	LastPage    int `json:"last_page"`
}
