package models

// RequestData represents the structure of the JSON object
type RequestData struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	HTTPStatus  int    `json:"httpStatus"`
}
