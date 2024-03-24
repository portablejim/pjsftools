package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func getFieldPermsForField(fieldName string, sfInfo TokenInfo, getWrapper AuthHttpGetter) ([]sfFieldPermissons, error) {
	queryStringRaw := fmt.Sprintf(`SELECT Id, ParentId, Parent.Profile.Name, Field, PermissionsEdit, PermissionsRead
		FROM FieldPermissions
		WHERE Parent.Type = 'Profile'
		AND Field='%s'
		ORDER BY Parent.Profile.Name`, fieldName)
	queryString := tidyForQueryParam(queryStringRaw)
	targetUrl := fmt.Sprintf("%s/services/data/v%s/query?q=%s", sfInfo.instanceUrl, sfInfo.version, queryString)
	fetchResp, fetchRespErr := getWrapper.AuthedGet(targetUrl, sfInfo.accessToken)
	if fetchRespErr != nil {
		return []sfFieldPermissons{}, fmt.Errorf("error executing GET Request")
	}

	fetchRespStr, fetchRespStrErr := io.ReadAll(fetchResp.Body)
	if fetchRespStrErr != nil {
		return []sfFieldPermissons{}, fmt.Errorf("error reading GET Request")
	}
	var currentPermsResult sfQueryResult[sfFieldPermissons]
	currentPermsParseErr := json.Unmarshal(fetchRespStr, &currentPermsResult)
	if currentPermsParseErr != nil {
		return []sfFieldPermissons{}, fmt.Errorf("error parsing GET Request. Error: %w", currentPermsParseErr)
	}

	if !currentPermsResult.Done {
		return []sfFieldPermissons{}, fmt.Errorf("incomplete result set. Extra result parsing not implemented")
	}

	return currentPermsResult.Records, nil
}
