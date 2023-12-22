package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicholasboyce/shufflify-cli/main"
)

func TestFetchWebAPI(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)

		switch r.URL.Path {
		case "/v1/me": //get profile info
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"display_name": "niko",
				"id": "417",
				"uri": "spotify:417",
				"product": "free",
				"images": [
					{
						"url": "https://iida.com/images/1234567"
					}
				]
			}`)
		case "/v1/me/playlists": //get current user's playlists
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"items": [
					{
						"name": "Groovy!",
						"uri": "spotify:572",
						"public": false,
						"images": [
							{
								"url": "https://yamashitarocks.jp"
							}
						],
						"id": "572"
					}
				]
			}`)
		case "/v1/playlists/playlist_id/tracks":
			if r.Method == http.MethodGet { // get items from playlist
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"items": [
						{
							"name": "Groovy!",
							"uri": "spotify:572",
							"public": false,
							"images": [
								{
									"url": "https://yamashitarocks.jp"
								}
							],
							"id": "572",
						}
					]
				}`))
			} else if r.Method == http.MethodPost { //add to items playlist
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{
					"snapshot_id": "abc"
				}`))
			}
		case "/v1/users/user_id/playlists": //create a playlist for user
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{
				"external_urls": {
					"spotify": "shufflify.app"
				}
			}`))
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))

	defer server.Close()

	client := server.Client()

	t.Run("Checking GET profile info call", func(t *testing.T) {
		type ImageData struct {
			URL string
		}

		type ProfileInfo struct {
			Display_name string      `json:"display_name"`
			ID           int         `json:"id,string"`
			URI          string      `json:"uri"`
			Product      string      `json:"product"`
			Images       []ImageData `json:"images"`
		}

		profileInfo := ProfileInfo{}

		_, err := main.FetchWebAPI("GET", fmt.Sprintf("%s/v1/me", server.URL), nil, &profileInfo, client)
		if err != nil {
			t.Error(err)
		}

		if profileInfo.Display_name != "niko" {
			t.Errorf("Expected 'niko', got %q\n", profileInfo.Display_name)
		}
		if profileInfo.ID != 417 {
			t.Errorf("Expected '417', got %v\n", profileInfo.ID)
		}
		if profileInfo.Images[0].URL != "https://iida.com/images/1234567" {
			t.Errorf("Expected 'https://iida.com/images/1234567', got %s\n", profileInfo.Images[0].URL)
		}
		if profileInfo.URI != "spotify:417" {
			t.Errorf("Expected 'spotify:417', got %q\n", profileInfo.URI)
		}
		if profileInfo.Product != "free" {
			t.Errorf("Expected 'free', got %q\n", profileInfo.Product)
		}
		fmt.Println(profileInfo)
		fmt.Println("")
	})

	t.Run("Checking GET user's playlists", func(t *testing.T) {
		type ImageData struct {
			URL string
		}

		type Playlist struct {
			Name   string      `json:"name"`
			URI    string      `json:"uri"`
			Public bool        `json:"public"`
			Images []ImageData `json:"images"`
			ID     int         `json:"id,string"`
		}

		type PlaylistItems struct {
			Items []Playlist `json:"items"`
		}

		playlistItems := PlaylistItems{}

		_, err := main.FetchWebAPI("GET", fmt.Sprintf("%s/v1/me/playlists", server.URL), nil, &playlistItems, client)
		if err != nil {
			t.Error(err)
		}

		fmt.Println(playlistItems)
		fmt.Println("")

		if playlistItems.Items[0].Name != "Groovy!" {
			t.Errorf("Expected 'Groovy!', got %q\n", playlistItems.Items[0].Name)
		}
		if playlistItems.Items[0].URI != "spotify:572" {
			t.Errorf("Expected 'spotify:572', got %q\n", playlistItems.Items[0].URI)
		}
		if playlistItems.Items[0].Public != false {
			t.Errorf("Expected %t, got %v", false, playlistItems.Items[0].Public)
		}
		if playlistItems.Items[0].ID != 572 {
			t.Errorf("Expected '572', got %v", playlistItems.Items[0].ID)
		}
		if playlistItems.Items[0].Images[0].URL != "https://yamashitarocks.jp" {
			t.Errorf("Expected 'https://yamashitarocks.jp', got %v", playlistItems.Items[0].Images[0].URL)
		}
	})

	t.Run("Checking POST items to user's playlist", func(t *testing.T) {
		var response struct {
			Snapshot_id string `json:"snapshot_id"`
		}

		main.FetchWebAPI("POST", fmt.Sprintf("%s/v1/playlists/playlist_id/tracks", server.URL), nil, &response, client)

		fmt.Println(response)
		fmt.Println("")

		if response.Snapshot_id != "abc" {
			t.Errorf("Expected 'abc', got %v", response.Snapshot_id)
		}
	})

	t.Run("Checking POST playlist to user's account", func(t *testing.T) {
		type ExternalURLs struct {
			Spotify string `json:"spotify"`
		}

		type Response struct {
			External_urls ExternalURLs `json:"external_urls"`
		}

		var response Response
		{
		}

		main.FetchWebAPI("POST", fmt.Sprintf("%s/v1/users/user_id/playlists", server.URL), nil, &response, client)

		fmt.Println(response)
		fmt.Println("")

		if response.External_urls.Spotify != "shufflify.app" {
			t.Errorf("Expected 'shufflify.app', got %v", response.External_urls.Spotify)
		}

	})

}
