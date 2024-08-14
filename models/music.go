package models

type Music struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Location string `json:"location"`
	Year     int32  `json:"year"`
}
