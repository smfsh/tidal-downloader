package main

import (
	"fmt"
)

const TIDAL_URL_BASE string = "https://api.tidalhifi.com/v1/"

func main() {
	fmt.Println("Starting Tidal Downloader")
	//setConfig()
	c := NewTidalConfig("LOSSLESS", "username", "password")

	login(c)
}
