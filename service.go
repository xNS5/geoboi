package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	ipapiClient := http.Client{}
	req, err := http.NewRequest("GET", "https://ipapi.co/json/", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "ipapi.co/#go-v1.3")
	resp, err := ipapiClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
