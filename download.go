package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type tidalBase struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

type tidalTrack struct {
	tidalBase
	Duration    int `json:"duration"`
	TrackNumber int `json:"trackNumber"`
	Artist      struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"artist"`
	Album struct {
		Id    int    `json:"id"`
		Title string `json:"name"`
	} `json:"album"`
}

type tidalAlbum struct {
	tidalBase
	Tracks []int
}

type tidalStream struct {
	Url           string `json:"url"`
	EncryptionKey string `json:"encryptionKey"`
}

func get(url string, c *tidalConfig) []byte {
	var sid string

	if c.quality != "LOSSLESS" {
		sid = c.sessionid1
	} else {
		sid = c.sessionid2
	}

	req, err := http.NewRequest("GET", TIDAL_URL_BASE+url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Tidal-SessionId", sid)

	q := req.URL.Query()
	q.Add("sessionId", sid)
	q.Add("countryCode", "US")
	q.Add("soundQuality", c.quality)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}

// Track: 62437814
func getTrackInfo(id int, c *tidalConfig) tidalTrack {
	body := get("tracks/"+strconv.Itoa(id), c)
	var jsonMap tidalTrack
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		panic(err)
	}

	return jsonMap
}

// Album: 62437813
func getAlbumInfo(id int, c *tidalConfig) tidalAlbum {
	body := get("albums/"+strconv.Itoa(id), c)
	var jsonMap tidalAlbum
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		panic(err)
	}

	return jsonMap
}

// Artist: 5221673
func getArtistInfo(id int, c *tidalConfig) tidalBase {
	body := get("artists/"+strconv.Itoa(id), c)
	var jsonMap tidalBase
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		panic(err)
	}

	return jsonMap
}

func getStreamUrl(id int, c *tidalConfig) tidalStream {
	body := get("tracks/"+strconv.Itoa(id)+"/streamUrl", c)
	var jsonMap tidalStream
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		panic(err)
	}

	return jsonMap
}

func getStreamExtension(url string) string {
	if strings.Contains(url, ".flac?") {
		return ".flac"
	} else if strings.Contains(url, ".mp4?") {
		return ".mp4"
	} else {
		return ".m4a"
	}
}

func downloadFile(url string, output string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func downloadTrack(trackId int, albumId int, echo bool, c *tidalConfig) {
	track := getTrackInfo(trackId, c)
	album := getAlbumInfo(albumId, c)

	if echo {
		fmt.Println("Preparing Track Download:")
		fmt.Println("Track Artist  ", track.Artist.Name)
		fmt.Println("Album Title   ", album.Title)
		fmt.Println("Track Title   ", track.Title)
		fmt.Println("Track Number  ", track.TrackNumber)
		fmt.Println("Track ID      ", track.Id)
		fmt.Println("Duration      ", track.Duration)
	}

	stream := getStreamUrl(track.Id, c)

	trackNumber := strconv.Itoa(track.TrackNumber)
	if track.TrackNumber < 10 {
		trackNumber = "0" + strconv.Itoa(track.TrackNumber)
	}
	ext := getStreamExtension(stream.Url)
	outName := trackNumber + " - " + track.Title + ext
	tempName := outName + ".tmp"
	err := downloadFile(stream.Url, tempName)
	if err != nil {
		panic(err)
	}

	if stream.EncryptionKey != "" {
		key, iv := decryptToken(stream.EncryptionKey)
		decryptFile(tempName, outName, key, iv)
		err := os.Remove(tempName)
		if err != nil {
			panic(err)
		}
	} else {
		err := os.Rename(tempName, outName)
		if err != nil {
			panic(err)
		}
	}
}

func downloadAlbum() {

}

func downloadArtist() {

}
