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

func main() {
	timezoneDir := "/usr/share/zoneinfo"
    localtimePath := "/etc/localtime"

	sysTz, err := GetLocalIanaName()

	if err != nil {
		log.Fatalf("Error getting local iana name: %v", err)
		return
	}

	ipGeoTz, err := GetRemoteIanaName()

	if err != nil {
		log.Fatalf("Error getting remote iana name: %v", err)
		return
	}
	
	// Checking to make sure that someone isn't trying to pass malicious code in the ipGeoTz string
	if isValid, err := ValidateIanaName(ipGeoTz); err != nil {
		log.Fatalf("Error validating IANA name: %v", err)
	} else if !isValid {
		log.Fatalf("Timezone is not valid: %s", ipGeoTz)
	} else if (ipGeoTz != sysTz) {
		timezonePath := filepath.Join(timezoneDir, ipGeoTz)

		if _, err := os.Stat(timezonePath); os.IsNotExist(err) {
			log.Fatalf("Timezone path does not exist: %v", err)
			return
		} else {
			log.Printf("Valid Timezone: %s", ipGeoTz)
		}
		
		if err := os.Remove(localtimePath); err != nil {
			log.Fatalf("Error removing current localtime: %v", err)
			return
		} else if err := os.Symlink(timezonePath, localtimePath); err != nil {
			log.Fatalf("Error creating symlink: %v", err)
			return
		} else {
			log.Printf("Timezone successfully changed to: %s", ipGeoTz)
		}
	} else {
		log.Println("Time zone unchanged")
	}
}

func ValidateIanaName(timezone string) (bool, error) {

	validTimezoneRegex := `^[A-Za-z]+(/[A-Za-z_-]+)+$`
    matched, err := regexp.MatchString(validTimezoneRegex, timezone)

	if err != nil {
        return false, fmt.Errorf("error compiling regex: %v", err)
    }

    if !matched {
        return false, fmt.Errorf("invalid timezone format: %s", timezone)
    }

	return true, nil
}

func GetRemoteIanaName() (string, error) {
	ipapiClient := http.Client{}

	if req, err := http.NewRequest("GET", "http://ip-api.com/json/", nil); err != nil {
		return "", err
	} else {
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
}

func GetLocalIanaName() (string, error) {
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

    if _, err = time.LoadLocation(timezone); err != nil {
		return "", err
	}
	
    return timezone, nil
}
