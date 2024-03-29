package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type sessionData struct {
	CountryCode string `json:"countryCode"`
	SessionId   string `json:"sessionId"`
	UserId      int    `json:"userId"`
}

func newSession() *tidalConfig {
	// Get username and password from the command line.
	username, password := getCredentials()

	// Prepare default configuration.
	c := newTidalConfig(false, "HI_RES", username, password)

	// Login with configuration.
	processLogin(c)

	return c
}

func getCredentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}

func processLogin(c *tidalConfig) {
	fmt.Println("\nAttempting to login as", c.username)
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

	resp, err := http.PostForm(TidalUrlBase+"login/username", params)
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
	validate, err := http.Get(TidalUrlBase + "users/" + strconv.Itoa(sessionData.UserId) + "?" + vp.Encode())
	if err != nil {
		panic(err)
	}
	if validate.StatusCode != 200 {
		panic("Tidal sessionId is invalid.")
	}

	return sessionData
}
