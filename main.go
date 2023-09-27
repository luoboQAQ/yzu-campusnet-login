package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func testConnection() bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", "http://111.13.141.31/generate_204", nil)
	if err != nil {
		panic(err)
	}
	req.Host = "connect.rom.miui.com"
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 204
}

func main() {
	LoadEnv()
	connected := false
	for {
		if !DEBUG && testConnection() {
			if !connected {
				log.Println("You have connected to the Internet")
				connected = true
			}
			time.Sleep(CHECK_INTERVAL)
			continue
		}
		connected = false

		log.Printf("Start login in %ds...\n", START_DELAY/time.Second)
		time.Sleep(START_DELAY)

		log.Printf("Username: %s, Password: %s, Service: %s\n", SSO_USERNAME, SSO_PASSWORD, CAMPUSNET_SERVICE)
		jar, _ := cookiejar.New(nil)
		client := &http.Client{
			Timeout: 5 * time.Second,
			Jar:     jar,
		}
		campusNet := NewCampusNet(client)
		sso := NewSSO(client)

		portalUrl, err := campusNet.GetPortalURL()
		if err != nil {
			log.Printf("Failed to get portal url: %v\n", err)
			continue
		}
		log.Printf("Portal url: %s\n", portalUrl)

		err = sso.Login(SSO_USERNAME, SSO_PASSWORD, portalUrl)
		if err != nil {
			log.Printf("Failed to login SSO: %v\n", err)
			continue
		}
		log.Println("Login SSO success")

		err = campusNet.LoginService(portalUrl, SSO_USERNAME, CAMPUSNET_SERVICE)
		if err != nil {
			log.Printf("Failed to login services: %v\n", err)
			continue
		}
		log.Println("Login service success")
	}
}
