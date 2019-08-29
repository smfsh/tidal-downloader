package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func get(url string, c *tidalConfig) map[string]interface{} {
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

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		panic(err)
	}

	return jsonMap
}

// Track: 62437814
func getTrack(id int, c *tidalConfig) map[string]interface{} {
	return get("tracks/"+strconv.Itoa(id), c)
}

// Album: 62437813
func getAlbum(id int, c *tidalConfig) map[string]interface{} {
	return get("albums/"+strconv.Itoa(id), c)
}

// Artist: 5221673
func getArtist(id int, c *tidalConfig) map[string]interface{} {
	return get("artists/"+strconv.Itoa(id), c)
}
