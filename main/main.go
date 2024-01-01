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

	names := flag.Args()

	userPlaylists := PlaylistItems{}
	userPlaylistNames := make(map[string]string)
	tracklist := []string{}

	//if flag.Args() >= 2, fetch user playlists. for every arg (playlist), if it's in user playlist set, add to tracklist.
	if len(names) >= 2 {

		profileInfo := ProfileInfo{}

		if _, err := FetchWebAPI("GET", "https://api.spotify.com/v1/me", nil, &profileInfo, client); err != nil {
			log.Fatal(err)
		}
		//fetch user playlists
		if _, err := FetchWebAPI("GET", "https://api.spotify.com/v1/me/playlists", nil, &userPlaylists, client); err != nil {
			log.Fatal(err)
		}

		for _, userPlaylist := range userPlaylists.Items {
			userPlaylistNames[userPlaylist.Name] = userPlaylist.ID
		}

		for _, name := range names {
			if id, present := userPlaylistNames[name]; present {
				playlistTracks := TrackItems{}
				FetchWebAPI("GET", fmt.Sprintf("https://api.spotify.com/v1/playlists/%v/tracks", id), nil, &playlistTracks, client)
				for _, track := range playlistTracks.Items {
					tracklist = append(tracklist, track.ID)
				}
			}
		}

		//shuffle tracklist. then ask for title and description of playlist. post to user account.
		Shuffle(tracklist)

		fmt.Println("Name of new playlist: ")
		var playlistTitle string
		if _, err := fmt.Scan(&playlistTitle); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v's description: \n", playlistTitle)
		var playlistDescription string
		if _, err := fmt.Scan(&playlistDescription); err != nil {
			log.Fatal(err)
		}

		playlistData := map[string]string{
			"name":        playlistTitle,
			"description": playlistDescription,
		}

		//create user playlist
		playlistPostResponse := struct {
			ID string `json:"id"`
		}{}

		if status, _ := FetchWebAPI("POST", fmt.Sprintf("https://api.spotify.com/v1/users/%v/playlists", profileInfo.ID), playlistData, &playlistPostResponse, client); status != http.StatusCreated {
			fmt.Errorf("Got status %s", status)
		}

		//add tracks to user playlist
		FetchWebAPI("POST", fmt.Sprintf("https://api.spotify.com/v1/playlists/playlist_id/tracks"), tracklist, nil, client)
	}

	//post url to playlist for user to click on/copy.
	//if premium account, ask if they want to queue tracks. if yes, for each track in tracklist post to queue.
	//when done, print thank you message and say 'all done!'

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
