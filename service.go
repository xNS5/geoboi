package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	timezoneDir := "/usr/share/zoneinfo"
    localtimePath := "/etc/localtime"

	sysTz, err := GetLocalIanaName()

	if err != nil {
		log.Fatal(err)
		return
	}

	ipGeoTz := GetRemoteIanaName()
	
	// Checking to make sure that someone isn't trying to pass malicious code in the ipGeoTz string
	isValid, err := ValidateIanaName(ipGeoTz)

	if !isValid {
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
		log.Println("Time zones unchanged")
	}
}

func main() {
	isOnline := IsOnline()

	if !isOnline {
		for i := 0; i < 30; i++ {
			log.Println("Unable to connect to internet, trying again...")
			time.Sleep(2 * time.Second)
			if IsOnline() == true {
				execTzChange()
				return
			}
		}
		log.Fatalln("Unable to connect to the internet")
	} else {
		execTzChange()
	}
}
