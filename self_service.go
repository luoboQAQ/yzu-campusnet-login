package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const SERVICE_HOST = "10.245.1.113:8080/selfservice/module"

type SelfService struct {
	client *http.Client
}

func NewSelfService(client *http.Client) *SelfService {
	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	return &SelfService{client}
}

func (s *SelfService) Login(username, password string) error {
	paramsData := url.Values{
		"name":     {username},
		"password": {password},
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/scgroup/web/login_judge.jsf?mobileslef=true", SERVICE_HOST), strings.NewReader(paramsData.Encode()))
	if err != nil {
		return fmt.Errorf("error creat request: %v", err)
	}
	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error self service login: %v", err)
	}
	defer resp.Body.Close()
	html, err := RespToString(resp)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`errorMsg=(.*?)&`)
	errorCode := re.FindStringSubmatch(html)
	if errorCode != nil {
		return fmt.Errorf("Login error:%s", errorCode)
	}
	return nil
}

func (s *SelfService) GetOnlines() ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/webcontent/web/onlinedevice_list.jsf", SERVICE_HOST), nil)
	if err != nil {
		return nil, fmt.Errorf("error creat request: %v", err)
	}
	req.Header.Add("User-Agent", USER_AGENT)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error get onlines: %v", err)
	}
	defer resp.Body.Close()
	html, err := RespToString(resp)
	if err != nil {
		return nil, err
	}

	if strings.Contains(html, "操作执行失败") {
		return nil, fmt.Errorf("get onlines error: 远程服务器错误")
	}
	re := regexp.MustCompile(`<input.*?id="userIp.*?value="(.*?)".*?>`)
	ips := re.FindAllStringSubmatch(html, -1)
	if ips == nil {
		return nil, fmt.Errorf("get onlines error: cannot find ips")
	}
	var result []string
	for _, ip := range ips {
		result = append(result, ip[1])
	}
	return result, nil
}

func (s *SelfService) Logout(name, ip string) error {
	paramsData := url.Values{
		"key": {name + ":" + ip},
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/userself/web/userself_ajax.jsf?methodName=indexBean.kickUserBySelfForAjax", SERVICE_HOST), strings.NewReader(paramsData.Encode()))
	if err != nil {
		return fmt.Errorf("error creat request: %v", err)
	}
	req.Header.Add("User-Agent", USER_AGENT)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error self service logout: %v", err)
	}
	defer resp.Body.Close()
	html, err := RespToString(resp)
	if err != nil {
		return err
	}

	if strings.Contains(html, "操作执行失败") {
		return fmt.Errorf("logout error: 远程服务器错误")
	}
	return nil
}
