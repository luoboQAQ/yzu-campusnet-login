package main

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var ERROR_MSG_CODE = map[string]string{
	"1320009": "验证码有误，请确认后重新输入",
	"1320010": "手机号未绑定用户，请使用其他方式登录",
	"1320011": "用户未绑定过手机号，请使用其他方式登录",
	"1320012": "短信发送服务存在问题，请稍后再试",
	"1320013": "用户名有误，请确认后重新输入",
	"1320033": "手机验证码已过期，请重新获取",
	"1030028": "账号被锁定",
	"1030048": "秒内不能重复获取验证码",
	"9280078": "验证码有误，请确认后重新输入",
	"9280081": "验证码已过期，请重新获取验证码",
	"2040001": "短信发送失败",
	"2040002": "验证码有误，请确认后重新输入",
	"2040003": "验证码有误，请确认后重新输入",
	"2040004": "手机号码不符合规范",
	"1030027": "用户名或密码错误，请确认后重新输入",
	"1030031": "用户名或密码错误，请确认后重新输入",
	"1410041": "当前用户名已失效",
	"1410040": "当前用户名已失效",
	"1320007": "验证码有误，请确认后重新输入",
}

const SSO_IP = "58.192.134.14"
const SSO_HOST = "sso.yzu.edu.cn"

// PKCS5Padding 对数据进行PKCS5填充
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func Encrypt(crypto, password string) string {
	key, err := base64.StdEncoding.DecodeString(crypto)
	if err != nil {
		panic(err)
	}

	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}

	data := []byte(password)

	bs := block.BlockSize()
	//对明文数据进行补码
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		panic("Need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//对明文按照blocksize进行分块加密
		//必要时可以使用go关键字进行并行加密
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	encodedResult := base64.StdEncoding.EncodeToString(out)
	return encodedResult
}

func SearchParams(id, htmlText string) string {
	re := regexp.MustCompile(fmt.Sprintf(`<p id="%s">(.*?)</p>`, id))
	match := re.FindStringSubmatch(htmlText)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

type SSO struct {
	client *http.Client
}

func NewSSO(client *http.Client) *SSO {
	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	return &SSO{client}
}

func (s *SSO) GetLoginParams(service string) (map[string]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/login", SSO_HOST), nil)
	if err != nil {
		return nil, fmt.Errorf("error logging in to service: %v", err)
	}
	paramsData := url.Values{
		"service": {service},
	}
	req.URL.RawQuery = paramsData.Encode()
	req.Header.Add("User-Agent", USER_AGENT)
	req.Host = SSO_HOST

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error logging in to service: %s", resp.Status)
	}
	html, err := RespToString(resp)
	if err != nil {
		return nil, err
	}
	params := make(map[string]string)
	params["execution"] = SearchParams("login-page-flowkey", html)
	params["croypto"] = SearchParams("login-croypto", html)
	return params, nil
}

func (s *SSO) Login(username, password, service string) error {
	params, err := s.GetLoginParams(service)
	if err != nil {
		return err
	}
	paramsData := url.Values{
		"username":    {username},
		"type":        {"UsernamePassword"},
		"_eventId":    {"submit"},
		"geolocation": {""},
		"execution":   {params["execution"]},
		"croypto":     {params["croypto"]},
		"password":    {Encrypt(params["croypto"], password)},
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/login", SSO_HOST), strings.NewReader(paramsData.Encode()))
	if err != nil {
		return fmt.Errorf("error sso login: %v", err)
	}
	req.Header.Add("User-Agent", USER_AGENT)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sso login: %v", err)
	}
	defer resp.Body.Close()
	html, err := RespToString(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized {

		re := regexp.MustCompile(`<div[^>]*?id="login-error-msg"[^>]*?>\s*<span>(.*?)</span>`)
		errorCode := re.FindStringSubmatch(html)[1]
		if _, ok := ERROR_MSG_CODE[errorCode]; !ok {
			return fmt.Errorf("unknown error code: %s", errorCode)
		}
		return errors.New(ERROR_MSG_CODE[errorCode])
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error sso login: %s", resp.Status)
	}
	return nil
}
