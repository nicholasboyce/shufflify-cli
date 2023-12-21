package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func FetchWebAPI(method string, resource string, body map[string]string, target interface{}, accessToken string) (string, error) {

	data := url.Values{}
	for key, value := range body {
		data.Set(key, value)
	}

	request, err := http.NewRequest(method, resource, strings.NewReader(data.Encode()))

	if err != nil || request == nil {
		fmt.Println(err)
		return "", err
	}

	request.Close = true

	request.Header.Add("Authorization", fmt.Sprintf("Bearer: %v", accessToken))

	// fmt.Print(request)
	client := &http.Client{
		Timeout: 30 * time.Second,
	} //TODO: Understand why I need to do this when client should be defined during login process

	response, fetchErr := client.Do(request)
	if fetchErr != nil {
		return "fetch", fetchErr
	}
	// fmt.Println(response)
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return response.Status, json.Unmarshal(content, &target)
}
