package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/browser"
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
		RedirectURL: "http://localhost:5173/oauth/callback",
	}

	// use PKCE to protect against CSRF attacks
	// https://www.ietf.org/archive/id/draft-ietf-oauth-security-topics-22.html#name-countermeasures-6
	verifier := oauth2.GenerateVerifier()

	//Create server which will collect the code for you
	codeChan := make(chan string)

	server := &http.Server{Addr: ":5173"}

	http.HandleFunc("/oauth/callback", handleOauthCallback(ctx, conf, codeChan))

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	fmt.Printf("Your browser has been opened to visit: %v\n", url)

	if err := browser.OpenURL(url); err != nil {
		panic(fmt.Errorf("failed to open browser for authentication %s", err.Error()))
	}

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.

	code := <-codeChan

	tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	if err := saveTokenAndConfig(tok, conf, path); err != nil {
		log.Fatal(err)
	}
	return client
}

func handleOauthCallback(ctx context.Context, config *oauth2.Config, codeChan chan string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParts, _ := url.ParseQuery(r.URL.RawQuery)

		// Use the authorization code that is pushed to the redirect URL.
		code := queryParts["code"][0]
		log.Printf("code: %s\n", code)

		// write the authorization code to the channel
		codeChan <- code

		msg := "<p><strong>Authentication successful</strong>. You may now close this tab.</p>"
		// send a success message to the browser
		fmt.Fprint(w, msg)
	}
}

func saveTokenAndConfig(token *oauth2.Token, conf *oauth2.Config, path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to create token file: %v", err)
	}
	defer file.Close()

	return EncodeTokenAndConfig(file, token, conf)
}

func EncodeTokenAndConfig(file io.Writer, token interface{}, conf interface{}) error {
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
	token := &oauth2.Token{}
	conf := &oauth2.Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open file at %v: %v", path, err)
	}
	defer file.Close()

	if err := DecodeTokenAndConfig(file, token, conf); err != nil {
		return nil, nil, err
	}

	return token, conf, nil
}

func DecodeTokenAndConfig(file io.Reader, token interface{}, conf interface{}) error {
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
	if err := saveTokenAndConfig(tok, conf, path); err != nil {
		log.Fatal(err)
	}
	return client
}
