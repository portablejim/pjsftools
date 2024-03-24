package main

import (
	"encoding/json"
	"net/url"
	"strings"
)

type TokenInfo struct {
	accessToken string
	version     string
	instanceUrl string
}

type sfOrgList struct {
	//status   int
	Result sfOrgListResult `json:"result"`
	//warnings []string
}

type sfOrgListResult struct {
	//id              string
	ApiVersion  string `json:"apiVersion"`
	AccessToken string `json:"accessToken"`
	InstanceUrl string `json:"instanceUrl"`
	//username        string
	//clientId        string
	ConnectedStatus string `json:"connectedStatus"`
	//sfdxAuthUrl     string
	//alias           string
}

type sfQueryResult[s any] struct {
	TotalSize int  `json:"totalSize"`
	Done      bool `json:"done"`
	Records   []s  `json:"records"`
}

type sfFieldPermissons struct {
	Id              string
	ParentId        string
	Parent          struct{ Profile struct{ Name string } }
	Field           string
	PermissionsRead bool
	PermissionsEdit bool
}

func getAccessToken(orgName string, runner CommandRunner) (TokenInfo, error) {
	cmdArgs := []string{"org", "display", "--verbose", "--json"}
	if len(orgName) > 0 {
		cmdArgs = append(cmdArgs, "-o")
		cmdArgs = append(cmdArgs, orgName)
	}

	cmdOutputBytes, cmdErr := runner.RunCommand("sf", cmdArgs)
	if cmdErr != nil {
		panic(cmdErr)
	}

	cmdOutput := sfOrgList{}
	json.Unmarshal(cmdOutputBytes, &cmdOutput)
	if cmdOutput.Result.ConnectedStatus != "Connected" {
		panic("Not connected to org.")
	}

	return TokenInfo{
		accessToken: cmdOutput.Result.AccessToken,
		instanceUrl: cmdOutput.Result.InstanceUrl,
		version:     cmdOutput.Result.ApiVersion,
	}, nil
}

func tidyForQueryParam(inputString string) string {
	outputStr := ""
	for _, part := range strings.Split(inputString, "\n") {
		outputStr += " " + strings.Trim(part, " \t")
	}

	return url.QueryEscape(outputStr[1:])
}
