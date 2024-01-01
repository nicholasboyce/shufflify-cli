package main

type ImageData struct {
	URL string
}

type Playlist struct {
	Name   string      `json:"name"`
	URI    string      `json:"uri"`
	Public bool        `json:"public"`
	Images []ImageData `json:"images"`
	ID     string      `json:"id"`
}

type PlaylistItems struct {
	Items []Playlist `json:"items"`
}

type ProfileInfo struct {
	Display_name string      `json:"display_name"`
	ID           string      `json:"id"`
	URI          string      `json:"uri"`
	Product      string      `json:"product"`
	Images       []ImageData `json:"images"`
}

type Track struct {
	ID string `json:"id"`
}

type TrackItems struct {
	Items []Track `json:"items"`
}
