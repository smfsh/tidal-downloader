package main

import (
	"fmt"
)

const TIDAL_URL_BASE string = "https://api.tidalhifi.com/v1/"

func main() {
	fmt.Println("Starting Tidal Downloader")

	// Get username and password from the command line.
	username, password := getCredentials()

	// Prepare default configuration.
	c := newTidalConfig("HI_RES", username, password)

	// Login with configuration.
	login(c)

	downloadTrack(80185438, 80185437, true, c)
}
