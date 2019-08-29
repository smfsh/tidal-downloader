package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type tidalConfig struct {
	sessionid1  string
	sessionid2  string
	uniqueKey   string
	countrycode string
	quality     string
	username    string
	password    string
	userid      int
}

type sessionData struct {
	CountryCode string `json:"countryCode"`
	SessionId   string `json:"sessionId"`
	UserId      int    `json:"userId"`
}

const TIDAL_URL_BASE string = "https://api.tidalhifi.com/v1/"

var config tidalConfig

func main() {
	fmt.Println("Starting Tidal Downloader")
	setConfig()

	login(config.username, config.password)
}

func setConfig() {
	config.uniqueKey = strings.Replace(uuid.New().String(), "-", "", -1)[16:]
	config.quality = "LOSSLESS"
	config.username = "username"
	config.password = "password"
}

func login(username string, password string) {
	fmt.Println("Attempting to login as", username)
	session1 := getSession(username, password, 1)
	session2 := getSession(username, password, 2)

	config.countrycode = session1.CountryCode
	config.userid = session1.UserId
	config.sessionid1 = session1.SessionId

	config.sessionid2 = session2.SessionId

	fmt.Println(config)
}

func getSession(username string, password string, session int) sessionData {
	var token string
	switch session {
	case 1:
		token = "4zx46pyr9o8qZNRw"
	case 2:
		token = "kgsOOmYk3zShYrNP"
	}

	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)
	params.Add("token", token)
	params.Add("clientUniqueKey", config.uniqueKey)
	params.Add("version", "1.9.1")

	resp, err := http.PostForm(TIDAL_URL_BASE+"login/username", params)
	if err != nil {
		panic("Unable to authenticate to Tidal.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Cannot read response body.")
	}
	sessionData := sessionData{}
	err = json.Unmarshal(body, &sessionData)
	if err != nil {
		panic(err)
	}

	vp := url.Values{}
	vp.Add("sessionId", sessionData.SessionId)
	validate, err := http.Get(TIDAL_URL_BASE + "users/" + strconv.Itoa(sessionData.UserId) + "?" + vp.Encode())
	if err != nil {
		panic(err)
	}
	if validate.StatusCode != 200 {
		panic("Tidal sessionId is invalid.")
	}

	return sessionData
}
