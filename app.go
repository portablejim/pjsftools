package main

import (
	"net/http"
	"os/exec"
	"strings"
)

type Command struct {
}

type CommandRunner interface {
	RunCommand(cmd string, cmdArgs []string) ([]byte, error)
}

func (c Command) RunCommand(cmd string, cmdArgs []string) ([]byte, error) {
	return exec.Command(cmd, cmdArgs...).Output()
}

type AuthHttpGetter interface {
	AuthedGet(url string, accessToken string) (http.Response, error)
}

func (c Command) AuthedGet(url string, accessToken string) (http.Response, error) {
	fetchReq, fetchReqErr := http.NewRequest("GET", url, nil)
	if fetchReqErr != nil {
		return http.Response{}, fetchReqErr
	}
	fetchReq.Header.Add("Authorization", "Bearer "+accessToken)
	fetchResp, fetchRespErr := http.DefaultClient.Do(fetchReq)

	return *fetchResp, fetchRespErr
}

type AuthHttpPatcher interface {
	AuthedPatch(url string, bodyStr string, contentType string, accessToken string) (http.Response, error)
}

func (c Command) AuthedPatch(url string, bodyStr string, contentType string, accessToken string) (http.Response, error) {
	fetchReq, fetchReqErr := http.NewRequest("POST", url, strings.NewReader(bodyStr))
	if fetchReqErr != nil {
		return http.Response{}, fetchReqErr
	}
	if len(contentType) == 0 {
		contentType = "application/json"
	}
	fetchReq.Header.Add("Authorization", "Bearer "+accessToken)
	fetchReq.Header.Add("Content-Type", contentType)
	fetchResp, fetchRespErr := http.DefaultClient.Do(fetchReq)

	return *fetchResp, fetchRespErr
}
