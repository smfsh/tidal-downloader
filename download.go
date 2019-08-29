package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type tidalObject struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Duration    int    `json:"duration"`
	TrackNumber int    `json:"trackNumber"`
	Url         string `json:"url"`
}

func get(url string, c *tidalConfig) tidalObject {
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
	//fmt.Printf("%s\n", body)

	var jsonMap tidalObject
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		panic(err)
	}

	return jsonMap
}

// Track: 62437814
func getTrack(id int, c *tidalConfig) tidalObject {
	return get("tracks/"+strconv.Itoa(id), c)
}

// Album: 62437813
func getAlbum(id int, c *tidalConfig) tidalObject {
	return get("albums/"+strconv.Itoa(id), c)
}

// Artist: 5221673
func getArtist(id int, c *tidalConfig) tidalObject {
	return get("artists/"+strconv.Itoa(id), c)
}

func getStreamUrl(id int, c *tidalConfig) tidalObject {
	return get("tracks/"+strconv.Itoa(id)+"/streamUrl", c)
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
