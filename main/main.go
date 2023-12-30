package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	var client *http.Client
	var path string

	flag.StringVar(&path, "filepath", "", "Path to the config file for auth detail storage")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if path == "" {
		path = os.Getenv("PATH_TO_AUTH_DETAILS")
	}

	if path == "" {
		path = "~/dev/.config/shufflify"
	}

	if !pathExists(path) {
		client = LoginProcess(path)
	} else {
		client = createNewClient(path)
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

func pathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}
