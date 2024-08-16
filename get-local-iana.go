package main

import (
	"os"
	"strings"
	"errors"
	"fmt"
	"time"
)

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

    _, err = time.LoadLocation(timezone)
    if err != nil {
        return "", err
    }

    return timezone, nil
}