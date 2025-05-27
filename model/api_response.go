package model

type ApiResponse struct {
	Success      bool
	Data         interface{}
	ErrorMessage string
}
