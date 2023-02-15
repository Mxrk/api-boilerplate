package dto

type ErrorResponse struct {
	ErrorType  string `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
