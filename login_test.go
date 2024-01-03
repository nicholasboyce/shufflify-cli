package main

import (
	"bytes"
	"testing"

	"golang.org/x/oauth2"
)

func TestTokenAndConfig(t *testing.T) {

	var bytes bytes.Buffer

	//Order of tests is important!!! Must write to bytes before reading from bytes

	t.Run("Given a token, config, and path, save pointer structs to JSON file at path", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken: "bellybutton",
		}
		conf := &oauth2.Config{
			ClientID: "identification",
		}

		if err := EncodeTokenAndConfig(&bytes, token, conf); err != nil {
			t.Errorf("Got error: %v\n", err)
		}
		result := bytes.String()
		if result == "" {
			t.Error("Encoding failed - result is empty.\n")
		}
	})

	t.Run("Given a token, config, and file, decode the file into the appropriate structs", func(t *testing.T) {
		token := &oauth2.Token{}
		conf := &oauth2.Config{}

		if err := DecodeTokenAndConfig(&bytes, token, conf); err != nil {
			t.Errorf("Got error: %v\n", err)
		}
		if token.AccessToken != "bellybutton" {
			t.Errorf("Expected 'bellybutton', got %v\n", token.AccessToken)
		}
		if conf.ClientID != "identification" {
			t.Errorf("Expected 'identification', got %v\n", conf.ClientID)
		}
	})
}
