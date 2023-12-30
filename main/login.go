package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

func LoginProcess(path string) *http.Client {
	var clientID string

	fmt.Print("Please input Client ID: ")
	if _, err := fmt.Scan(&clientID); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID: clientID,
		Scopes:   []string{"user-read-private, user-read-playback-state, user-modify-playback-state, playlist-read-private, playlist-read-collaborative, playlist-modify-private, playlist-modify-public"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		RedirectURL: "http://localhost:5173",
	}

	// use PKCE to protect against CSRF attacks
	// https://www.ietf.org/archive/id/draft-ietf-oauth-security-topics-22.html#name-countermeasures-6
	verifier := oauth2.GenerateVerifier()

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.

	fmt.Println("Please input code parameter from redirect url: ")
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	saveTokenAndConfig(tok, conf, path)
	return client
}

func saveTokenAndConfig(token *oauth2.Token, conf *oauth2.Config, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to create token file: %v", err)
	}
	defer file.Close()

	return encodeTokenAndConfig(file, token, conf)
}

func encodeTokenAndConfig(file io.Writer, token interface{}, conf interface{}) error {
	// encode token as JSON and write to file
	encoder := json.NewEncoder(file)

	if err := encoder.Encode(token); err != nil {
		return fmt.Errorf("unable to write token to file: %v", err)
	}

	if err := encoder.Encode(conf); err != nil {
		return fmt.Errorf("unable to write config struct to file: %v", err)
	}

	return nil
}

func fetchTokenAndConfig(path string) (*oauth2.Token, *oauth2.Config, error) {
	var token *oauth2.Token
	var conf *oauth2.Config

	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open file at %v: %v", path, err)
	}
	defer file.Close()

	if err := decodeTokenAndConfig(file, token, conf); err != nil {
		return nil, nil, err
	}

	return token, conf, nil
}

func decodeTokenAndConfig(file io.Reader, token interface{}, conf interface{}) error {
	//decode from file and write to token and conf structs
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(token); err != nil {
		return fmt.Errorf("unable to decode token from file: %v", err)
	}

	if err := decoder.Decode(conf); err != nil {
		return fmt.Errorf("unable to decode config struct from file: %v", err)
	}

	return nil
}

func createNewClient(path string) *http.Client {
	ctx := context.Background()
	tok, conf, err := fetchTokenAndConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	client := conf.Client(ctx, tok)
	return client
}
