package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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

func ValidateIanaName(timezone string) (bool, error) {

	validTimezoneRegex := `^[A-Za-z]+(/[A-Za-z_-]+)+$`
    matched, err := regexp.MatchString(validTimezoneRegex, timezone)

	if err != nil {
        return false, fmt.Errorf("Error compiling regex: %v", err)
    }

    if !matched {
        return false, fmt.Errorf("Invalid timezone format: %s", timezone)
    }

	return true, nil
}

func execTzChange(){
	timezoneDir := "/usr/share/zoneinfo"
    // localtimePath := "/etc/localtime"
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
	tempTz := "America/Los_Angeles"
	// Checking to make sure that someone isn't trying to pass malicious code in the ipGeoTz string
	isValid, err := ValidateIanaName(tempTz)

	if !isValid {
		log.Printf("Timezone is not valid: %s", tempTz)
	} else if (ipGeoTz != sysTz) {
		timezonePath := filepath.Join(timezoneDir, ipGeoTz)
		fmt.Println(timezonePath)

		if _, err := os.Stat(tempTz); os.IsNotExist(err) {
			log.Fatalf("Timezone file does not exist: %s", tempTz)
		} else {
			log.Println("Timezone is valid")
		}
		// else if err := os.Remove(localtimePath); err != nil {
		// 	log.Fatal("Error removing current localtime: %v", err)
		// } else if err := os.Symlink(timezonePath, localtimePath); err != nil {
		// 	log.Fatal("Error creating symlink: %v", err)
		// }
		log.Println("Timezones do not match, changing....")
	}
}

func main() {
	isOnline := IsOnline()

	if isOnline {
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
