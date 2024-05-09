package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Info struct {
	Code           int    `json:"code"`
	Msg            string `json:"msg"`
	OneDriveStatus int    `json:"oneDriveStatus"`
	TodoCount      string `json:"todoCount"`
	EcardBalance   string `json:"ecardBalance"`
	BorrowCount    string `json:"borrowCount"`
	NewMailCount   string `json:"newMailCount"`
	SamInfo        struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			FreeTotal       string    `json:"freeTotal"`
			FreeRemain      string    `json:"freeRemain"`
			UserID          string    `json:"userId"`
			UserName        string    `json:"userName"`
			AccountFee      float64   `json:"accountFee"`
			PreAccountFee   float64   `json:"preAccountFee"`
			PeriodStartTime time.Time `json:"periodStartTime"`
			PeriodEndTime   time.Time `json:"periodEndTime"`
			RemainDays      int       `json:"remainDays"`
		} `json:"data"`
	} `json:"samInfo"`
}

type QueryInfo struct {
	client *http.Client
}

func NewQueryInfo(client *http.Client) *QueryInfo {
	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	return &QueryInfo{client}
}

func (q *QueryInfo) GetToken() (string, error) {
	req, err := http.NewRequest("GET", "https://i.yzu.edu.cn/Student", nil)
	if err != nil {
		return "", fmt.Errorf("error get Request Token: %v", err)
	}
	req.Header.Add("User-Agent", USER_AGENT)

	resp, err := q.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error get Request Token: %s", resp.Status)
	}
	html, err := RespToString(resp)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("<input name=\"__RequestVerificationToken\" type=\"hidden\" value=\"(.*?)\" />")
	match := re.FindStringSubmatch(html)
	if len(match) >= 2 {
		return match[1], nil
	}
	return "", nil
}

func (q *QueryInfo) GetInfoJson(token string) (Info, error) {
	var info Info
	paramsData := url.Values{
		"__RequestVerificationToken": {token},
	}
	req, err := http.NewRequest("POST", "https://i.yzu.edu.cn/Student/Home/GetEvents", strings.NewReader(paramsData.Encode()))
	if err != nil {
		return info, fmt.Errorf("error get Info Json: %v", err)
	}
	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := q.client.Do(req)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return info, fmt.Errorf("error get Request Token: %s", resp.Status)
	}
	jsonInput, err := RespToString(resp)
	if err != nil {
		return info, err
	}
	err = json.Unmarshal([]byte(jsonInput), &info)
	if err != nil {
		return info, fmt.Errorf("error get Info Json: %v", err)
	}
	return info, nil
}
