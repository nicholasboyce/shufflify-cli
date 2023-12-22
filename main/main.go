package main

import (
	"fmt"
	"net/http"
)

var loggedIn bool = false

func main() {

	var client *http.Client

	if !loggedIn {
		client = LoginProcess()
	}

	type ImageData struct {
		URL string
	}

	type ProfileInfo struct {
		Display_name string      `json:"display_name"`
		ID           string      `json:"id"`
		URI          string      `json:"uri"`
		Product      string      `json:"product"`
		Images       []ImageData `json:"images"`
	}

	profileInfo := ProfileInfo{}

	x, y := FetchWebAPI("GET", "https://api.spotify.com/v1/me", nil, &profileInfo, client)

	fmt.Println(x, y)
	fmt.Println(profileInfo)

}
