package main

import (
	"strings"

	"github.com/google/uuid"
)

type tidalConfig struct {
	sessionid1  string
	sessionid2  string
	uniqueKey   string
	countrycode string
	parallel    bool
	quality     string
	username    string
	password    string
	userid      int
}

func newTidalConfig(parallel bool, quality string, username string, password string) *tidalConfig {
	c := new(tidalConfig)

	c.uniqueKey = strings.Replace(uuid.New().String(), "-", "", -1)[16:]
	c.parallel = parallel
	c.quality = quality
	c.username = username
	c.password = password

	return c
}
