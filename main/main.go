package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	var client *http.Client
	var path string
	var logout bool

	flag.StringVar(&path, "filepath", "", "Path to the config file for auth detail storage")
	flag.BoolVar(&logout, "logout", false, "Logout command")
	flag.Parse()

	if logout {
		os.RemoveAll(os.Getenv("PATH_TO_AUTH_DETAILS"))
		os.Exit(0)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if path == "" {
		path = os.Getenv("PATH_TO_AUTH_DETAILS")
	}

	if path == "" {
		path = "config.json"
	}

	os.Setenv("PATH_TO_AUTH_DETAILS", path)

	if !pathValid(path) { // simulates login state
		client = LoginProcess(path)
	} else {
		client = createNewClient(path)
	}

	//if flag.Args() >= 2, fetch user playlists. for every arg (playlist), if it's in user playlist set, fetch tracks and add to tracklist.
	//shuffle tracklist. then ask for title and description of playlist. post to user account.
	//post url to playlist for user to click on/copy.
	//if premium account, ask if they want to queue tracks. if yes, for each track in tracklist post to queue.
	//when done, print thank you message and say 'all done!'

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

	//TODO: LOGOUT FUNCTION: deletes the config file at saved path in environment variable

}

func pathValid(path string) bool {
	now := time.Now()

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		tok, _, err := fetchTokenAndConfig(path)
		if err != nil {
			log.Fatal(err)
		}
		return now.Before(tok.Expiry)
	}
}
