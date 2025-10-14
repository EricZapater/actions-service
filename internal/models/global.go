package models

type ResponseMessage struct {
	Result  string      `json:"result" example:"success"`
	Message string      `json:"message" example:"Action completed successfully"`
	Content interface{} `json:"content,omitempty"`
}