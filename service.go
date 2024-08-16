package main

import (
	"log"
	"os"
	"path/filepath"
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
		log.Println("Time zone unchanged")
	}
}
