package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

const TidalUrlBase string = "https://api.tidalhifi.com/v1/"

func main() {
	fmt.Println("Starting Tidal Downloader")

	c := newSession()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a Tidal URL: ")
		tidalUrl, _ := reader.ReadString('\n')

		err := processUrl(tidalUrl, c)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func processUrl(tidalUrl string, c *tidalConfig) error {
	u, err := url.Parse(strings.TrimSpace(tidalUrl))
	if err != nil {
		panic(err)
	}
	id, err := strconv.Atoi(path.Base(u.Path))
	if err != nil {
		return errors.New("Input URL must end with ID...")
	}

	// https://tidal.com/browse/track/116415079
	if strings.Contains(u.Path, "track") {
		fmt.Println("Found a track URL, processing...")
		downloadTrack(id, tidalAlbum{}, true, c)
		return nil
	}
	// https://listen.tidal.com/album/116415070
	if strings.Contains(u.Path, "album") {
		fmt.Println("Found an album, processing...")
		downloadAlbum(id, c.parallel, c)
		return nil
	}
	// https://listen.tidal.com/artist/3850668
	if strings.Contains(u.Path, "artist") {
		fmt.Println("Found an artist URL, processing...")
		downloadArtist(id, c)
		return nil
	}

	return errors.New("Input URL must be a track, album, or artist...")
}
