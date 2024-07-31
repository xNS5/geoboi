package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func IsOnline() bool {
	timeout := time.Duration(5000 * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}
	_, err := client.Get("https://google.com")

	if err != nil {
		return false
	}

	return true
}

func execTzChange(){
	ipapiClient := http.Client{}

	req, err := http.NewRequest("GET", "https://ipapi.co/json/", nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("First Request", req)

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

func main() {
	isOnline := IsOnline()

	if isOnline == false {
		for i := 0; i < 100; i++ {
			time.Sleep(2 * time.Second)
			if IsOnline() == true {
				execTzChange()
				break
			}

		}
	} else {
		execTzChange()
	}
}
