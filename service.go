package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

func GetIanaName() (string, error) {
    linkPath := "/etc/localtime"
    targetPath, err := os.Readlink(linkPath)
    if err != nil {
        return "", err
    }

    tzParts := strings.Split(targetPath, "/")
    if len(tzParts) < 3 {
        return "", errors.New("invalid timezone format")
    }

    continent, country := tzParts[len(tzParts)-2], tzParts[len(tzParts)-1]
    timezone := fmt.Sprintf("%s/%s", continent, country)

    _, err = time.LoadLocation(timezone)
    if err != nil {
        return "", err
    }

    return timezone, nil
}

func execTzChange(){

	sysTz, err := GetIanaName()

	if err != nil {
		log.Fatal(err)
		return
	}

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
	var inputJson map[string]interface{}

	json.Unmarshal([]byte(string(body)), &inputJson)
	
	ipGeoTz := inputJson["timezone"].(string)

	if ipGeoTz != sysTz {
		os.Setenv("TZ", ipGeoTz)
	}

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
