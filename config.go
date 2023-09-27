package main

import (
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var (
	USER_AGENT        string
	SSO_USERNAME      string
	SSO_PASSWORD      string
	CAMPUSNET_SERVICE string
	CHECK_INTERVAL    time.Duration
	START_DELAY       time.Duration
	DEBUG             bool
)

func LoadEnv() {
	USER_AGENT = os.Getenv("USER_AGENT")
	if USER_AGENT == "" {
		USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0"
	}
	SSO_USERNAME = os.Getenv("SSO_USERNAME")
	if SSO_USERNAME == "" {
		panic("SSO_USERNAME is not set")
	}
	SSO_PASSWORD = os.Getenv("SSO_PASSWORD")
	if SSO_PASSWORD == "" {
		panic("SSO_PASSWORD is not set")
	}
	CAMPUSNET_SERVICE = os.Getenv("CAMPUSNET_SERVICE")
	if CAMPUSNET_SERVICE == "" {
		panic("CAMPUSNET_SERVICE is not set")
	}
	check_int := os.Getenv("CHECK_INTERVAL")
	if check_int == "" {
		CHECK_INTERVAL = 60 * time.Second
	} else {
		CHECK_INTERVAL, _ = time.ParseDuration(check_int)
	}
	start_delay := os.Getenv("START_DELAY")
	if start_delay == "" {
		START_DELAY = 5 * time.Second
	} else {
		START_DELAY, _ = time.ParseDuration(start_delay)
	}
	DEBUG = os.Getenv("DEBUG") == "true"
}
