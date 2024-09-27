package models

type AggregatedData struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Count       int32  `json:"count"`
}
