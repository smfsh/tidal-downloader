package main

import (
	"fmt"
)

const TIDAL_URL_BASE string = "https://api.tidalhifi.com/v1/"

func main() {
	fmt.Println("Starting Tidal Downloader")

	// Prepare default configuration.
	c := newTidalConfig("LOSSLESS", "username", "password")

	// Login with configuration.
	login(c)

	track := getTrack(62437814, c)
	//get("albums/62437813", c)
	album := getAlbum(62437813, c)

	//printMap("", track)
	//printMap("", album)

	fmt.Println("Preparing download:")
	fmt.Println("Album Title   ", album["title"])
	fmt.Println("Track Title   ", track["title"])
	fmt.Println("Duration      ", track["duration"])
	fmt.Println("Track Number  ", track["trackNumber"])
}

func printMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			printMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}
