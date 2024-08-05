package main

import (
	"fmt"
	"regexp"
)

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