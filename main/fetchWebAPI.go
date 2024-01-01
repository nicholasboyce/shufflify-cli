package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func FetchWebAPI(method string, resource string, body map[string]string, target interface{}, client *http.Client) (int, error) {

	data := url.Values{}
	for key, value := range body {
		data.Set(key, value)
	}

	request, err := http.NewRequest(method, resource, strings.NewReader(data.Encode()))

	if err != nil || request == nil {
		fmt.Println(err)
		return 0, err
	}

	request.Close = true

	response, fetchErr := client.Do(request)
	if fetchErr != nil {
		return 0, fetchErr
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	return response.StatusCode, json.Unmarshal(content, &target)
}
