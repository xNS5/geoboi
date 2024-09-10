package main

import (
	"encoding/json"
	"io"
	"net/http"
)


func GetRemoteIanaName() (string, error) {
	ipapiClient := http.Client{}

	req, err := http.NewRequest("GET", "http://ip-api.com/json/", nil)

	if err != nil {
		return "", err
	}

	resp, err := ipapiClient.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	
	var inputJson map[string]interface{}

	json.Unmarshal([]byte(string(body)), &inputJson)
	
	return inputJson["timezone"].(string), nil
}
