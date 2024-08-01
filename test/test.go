package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)


func changeTimezone(timezone string) error {
    timezoneDir := "/usr/share/zoneinfo"
    localtimePath := "/etc/localtime"

    // Regex pattern to validate timezone format
    validTimezonePattern := `^[A-Za-z]+(/[A-Za-z_-]+)+$`
    matched, err := regexp.MatchString(validTimezonePattern, timezone)
    if err != nil {
        return fmt.Errorf("error compiling regex: %v", err)
    }
    if !matched {
        return fmt.Errorf("invalid timezone format: %s", timezone)
    }

    // Construct the full path to the timezone file
    timezonePath := filepath.Join(timezoneDir, timezone)

    // Check if the timezone file exists
    if _, err := os.Stat(timezonePath); os.IsNotExist(err) {
        return fmt.Errorf("timezone file does not exist: %s", timezonePath)
    }

    // Remove the existing symlink or file
    if err := os.Remove(localtimePath); err != nil {
        return fmt.Errorf("error removing current localtime: %v", err)
    }

    // Create a new symlink to the desired timezone
    if err := os.Symlink(timezonePath, localtimePath); err != nil {
        return fmt.Errorf("error creating symlink: %v", err)
    }

    fmt.Println("Timezone successfully changed to", timezone)
    return nil
}



func main(){
    timezone := "America/New_York echo \"Hello, world!\"" // Malicious input

    err := changeTimezone(timezone)
    if err != nil {
        fmt.Println("Error:", err)
    }
}