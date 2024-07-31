# GeoBoi

This is a simple systemd service that automatically adjusts the system time zone.

## Why is this necessary?

In Pop!_OS 22.04, I found that the "automatic time zone" wasn't updating my system time zone automatically. As an experiment, I wanted to create a service using golang that determines my location based on my public IP address. 

## Why Golang?

Go seems to be growing in popularity, and I wanted to take it for a spin. I learn best by making something useful, so that's exactly what this is. 