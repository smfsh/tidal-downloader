package main

import (
	"github.com/google/uuid"
	"strings"
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

func NewTidalConfig(quality string, username string, password string) *tidalConfig {
	c := new(tidalConfig)

	c.uniqueKey = strings.Replace(uuid.New().String(), "-", "", -1)[16:]
	c.quality = quality
	c.username = username
	c.password = password

	return c
}
