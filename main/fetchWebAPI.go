package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func FetchWebAPI(method string, resource string, body map[string]any, target interface{}, client *http.Client) (int, error) {

	var data io.Reader

	if body != nil {
		info, err := json.Marshal(body)
		if err != nil {
			return 0, err
		}
		data = bytes.NewReader(info)
	} else {
		data = nil
	}

	request, err := http.NewRequest(method, resource, data)

	if err != nil || request == nil {
		log.Fatal(request)
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
