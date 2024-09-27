package models

// ParentData represents the structure of each object in the JSON array
type ParentData struct {
	Parent   string `json:"parent"`
	Children string `json:"children"`
}
