package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type sessionData struct {
	CountryCode string `json:"countryCode"`
	SessionId   string `json:"sessionId"`
	UserId      int    `json:"userId"`
}

func login(c *tidalConfig) {
	fmt.Println("Attempting to login as", c.username)
	session1 := getSession(c, 1)
	session2 := getSession(c, 2)

	c.countrycode = session1.CountryCode
	c.userid = session1.UserId
	c.sessionid1 = session1.SessionId

	c.sessionid2 = session2.SessionId

	fmt.Println("Settings:", c)
}

func getSession(c *tidalConfig, session int) sessionData {
	var token string
	switch session {
	case 1:
		token = "4zx46pyr9o8qZNRw"
	case 2:
		token = "kgsOOmYk3zShYrNP"
	}

	params := url.Values{}
	params.Add("username", c.username)
	params.Add("password", c.password)
	params.Add("token", token)
	params.Add("clientUniqueKey", c.uniqueKey)
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
