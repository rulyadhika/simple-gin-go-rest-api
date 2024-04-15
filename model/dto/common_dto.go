package dto

type ApiResponse struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}
