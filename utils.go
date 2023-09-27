package main

import (
	"io"
	"net/http"
)

func RespToString(resp *http.Response) (string, error) {
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
