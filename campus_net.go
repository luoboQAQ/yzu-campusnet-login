package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type CampusNet struct {
	client *http.Client
}

func NewCampusNet(client *http.Client) *CampusNet {
	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	return &CampusNet{client}
}

func (c *CampusNet) GetHostQuery(portalURL string) (string, string, error) {
	matchResult := regexp.MustCompile(`http://(.+?)/.*?\?(.+)`).FindStringSubmatch(portalURL)
	if matchResult == nil {
		return "", "", errors.New("invalid portal URL")
	}
	return matchResult[1], matchResult[2], nil
}

func (c *CampusNet) GetPortalURL() (string, error) {
	req, err := http.NewRequest("GET", "http://123.123.123.123", nil)
	if err != nil {
		return "", fmt.Errorf("cannot get portal URL: %v", err)
	}
	req.Header.Add("User-Agent", USER_AGENT)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot get portal URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cannot get portal URL: %s", resp.Status)
	}
	respText, err := RespToString(resp)
	if err != nil {
		return "", fmt.Errorf("cannot get portal URL: %v", err)
	}

	urlRegex := regexp.MustCompile(`href='(.*?)'`)
	urlMatch := urlRegex.FindStringSubmatch(respText)
	if urlMatch == nil {
		return "", errors.New("cannot get portal URL: cannot find URL")
	}
	return urlMatch[1], nil
}

func (c *CampusNet) LoginService(portalURL, userID, service string) error {
	host, queryString, err := c.GetHostQuery(portalURL)
	if err != nil {
		return err
	}

	formData := url.Values{
		"userId":          {userID},
		"flag":            {"casauthofservicecheck"},
		"service":         {url.QueryEscape(service)},
		"queryString":     {queryString},
		"operatorPwd":     {""},
		"operatorUserId":  {""},
		"passwordEncrypt": {"false"},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/eportal/InterFace.do?method=loginOfCas", host), strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("error logging in to service: %v", err)
	}
	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error logging in to service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error logging in to service: %s", resp.Status)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("error decoding JSON response: %v", err)
	}

	if data["result"].(string) != "success" {
		return errors.New(data["message"].(string))
	}

	return nil
}

func (c *CampusNet) OldLoginService(portalURL, userID, password, service string) error {
	host, queryString, err := c.GetHostQuery(portalURL)
	if err != nil {
		return err
	}

	formData := url.Values{
		"userId":          {userID},
		"password":        {password},
		"service":         {url.QueryEscape(service)},
		"queryString":     {queryString},
		"operatorPwd":     {""},
		"operatorUserId":  {""},
		"validcode":       {""},
		"passwordEncrypt": {"false"},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/eportal/InterFace.do?method=login", host), strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("error logging in to service: %v", err)
	}
	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error logging in to service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error logging in to service: %s", resp.Status)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("error decoding JSON response: %v", err)
	}

	if data["result"].(string) != "success" {
		return errors.New(data["message"].(string))
	}

	return nil
}