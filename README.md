# GeoBoi

This is a simple systemd service that automatically adjusts the system time zone.

## Why is this necessary?

In Pop!_OS 22.04, I found that the "automatic time zone" wasn't updating my system time zone automatically. As an experiment, I wanted to create a service using golang that determines my location based on my public IP address. 

## Why Golang?

Go seems to be growing in popularity, and I wanted to take it for a spin. I learn best by making something useful, so that's exactly what this is. 

# Build


## Prereqisites

All that's really needed is Golang installed on your system at `/usr/local/go/bin/go`, and all of your environment variables set correctly. Feel free to modify `build.sh` to your needs. Or don't, it's your life. 

## Installer Script

I created a simple shell script `build.sh` which handles all of the building + moving of files.

The script does the following:
1. Checks to see if the service file exists. If it doesn't, it clones the template file and changes the execuatble path to `$GOBIN/geoboi`
2. Builds geoboi and installs it to `$GOBIN`
3. Creates a symbolic link from `geoboi.service` to `/etc/systemd/system/`
4. Reloads systemctl with any changes in `geoboi.service`

## Run Instructions

I've found that Bash can sometimes not see the global environment variables, so run the script with the -E flag (e.g. `sudo -E bash build.sh`). I spent an uncomfortable amount of time figuring that out.
