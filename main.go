package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	var client *http.Client
	var path string
	var logout bool

	flag.StringVar(&path, "filepath", "", "Path to the config file for auth detail storage")
	flag.BoolVar(&logout, "logout", false, "Logout command")
	flag.Parse()

	if path == "" {
		path = os.Getenv("PATH_TO_CONFIG")
	}

	if path == "" {
		execPath, _ := os.Executable()
		parentFolder := filepath.Dir(execPath)
		path = fmt.Sprint(parentFolder, "/shufflify/config.json")
	}

	if logout {
		os.RemoveAll(path)
		os.Exit(0)
	}

	os.Setenv("PATH_TO_CONFIG", path)

	if !pathValid(path) { // login state
		client = LoginProcess(path)
	} else {
		client = createNewClient(path)
	}

	names := flag.Args()

	if len(names) >= 2 {

		userPlaylists := PlaylistItems{}
		userPlaylistNames := make(map[string]string)
		tracklist := []string{}

		profileInfo := ProfileInfo{}

		if _, err := FetchWebAPI("GET", "https://api.spotify.com/v1/me", nil, &profileInfo, client); err != nil {
			log.Fatal(err)
		}

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
				for _, trackObject := range playlistTracks.Items {
					tracklist = append(tracklist, trackObject.Track.URI)
				}
			}
		}

		Shuffle(tracklist)

		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Name of new playlist: ")
		scanner.Scan()
		playlistTitle := scanner.Text()

		fmt.Printf("%v's description: \n", playlistTitle)
		scanner.Scan()
		playlistDescription := scanner.Text()

		playlistData := map[string]interface{}{
			"name":        playlistTitle,
			"description": playlistDescription,
		}

		playlistPostResponse := Response{}

		if status, err := FetchWebAPI("POST", fmt.Sprintf("https://api.spotify.com/v1/users/%v/playlists", profileInfo.ID), playlistData, &playlistPostResponse, client); status != http.StatusCreated {
			fmt.Println(status)
			log.Fatal(err)
		}

		playlistTrackBodyData := map[string]interface{}{
			"uris": tracklist,
		}

		if status, err := FetchWebAPI("POST", fmt.Sprintf("https://api.spotify.com/v1/playlists/%v/tracks", playlistPostResponse.ID), playlistTrackBodyData, nil, client); status != http.StatusCreated {
			fmt.Println(status)
			log.Fatal(err)
		}

		fmt.Printf("Check out your new playlist at %v\n", playlistPostResponse.External_urls.Spotify)

		if profileInfo.Product == "premium" {
			fmt.Println("Would you like to queue the tracks in your new playlist? Answer yes if so. ")
			var shouldQueue string
			if _, err := fmt.Scan(&shouldQueue); err != nil {
				log.Fatal(err)
			}
			answer := strings.Trim(shouldQueue, " ")
			answer = strings.ToLower(answer)
			if answer == "yes" {
				for _, track := range tracklist {
					uri := url.Values{}
					uri.Add("uri", track)
					if status, err := FetchWebAPI("POST", fmt.Sprintf("https://api.spotify.com/v1/me/player/queue?%v", uri.Encode()), nil, nil, client); status != http.StatusNoContent {
						fmt.Println(status)
						log.Fatal(err)
					}
				}
			}
		}

		fmt.Println("Thank you for using Shufflify!")
	}

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
