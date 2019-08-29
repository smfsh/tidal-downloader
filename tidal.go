package main

import (
	"fmt"
)

const TIDAL_URL_BASE string = "https://api.tidalhifi.com/v1/"

func main() {
	fmt.Println("Starting Tidal Downloader")

	// Prepare default configuration.
	c := newTidalConfig("HI_RES", "username", "password")

	// Login with configuration.
	login(c)

	track := getTrack(80185439, c)
	//get("albums/62437813", c)
	album := getAlbum(80185437, c)

	//printMap("", track)
	//printMap("", album)

	fmt.Println("Preparing download:")
	fmt.Println("Album Title   ", album.Title)
	fmt.Println("Track Title   ", track.Title)
	fmt.Println("Track ID      ", track.Id)
	fmt.Println("Duration      ", track.Duration)
	fmt.Println("Track Number  ", track.TrackNumber)

	//var id int = track["id"]
	stream := getStreamUrl(track.Id, c)
	//fmt.Println(stream.Url)

	err := downloadFile(stream.Url, "file.flac")
	if err != nil {
		panic(err)
	}
}
