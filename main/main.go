package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

var loggedIn bool = false

func main() {

	var client *http.Client

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if !loggedIn {
		client = LoginProcess()
	} // else {
	// 	client = createNewClient()
	// }

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
