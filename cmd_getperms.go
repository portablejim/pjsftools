package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func GetPerms(args []string, runner CommandRunner, getWrapper AuthHttpGetter) error {

	fs := flag.NewFlagSet("getperms", flag.ExitOnError)
	orgName := fs.String("org", "", "sf org to use")
	fs.Parse(args)

	if len(*orgName) == 0 {
		os.Stderr.WriteString("Using default org\n")
	} else {
		os.Stderr.WriteString("Using org " + *orgName + "\n")
	}
	sfInfo, tokenError := getAccessToken(*orgName, runner)
	if tokenError != nil {
		return fmt.Errorf("unable to get access token. Error: %w", tokenError)
	}
	if len(sfInfo.accessToken) < 1 || len(sfInfo.instanceUrl) < 1 || len(sfInfo.version) < 1 {
		return fmt.Errorf("something is wrong. Unable to get valid info")
	}

	fieldName := fs.Arg(0)
	fmt.Fprintf(os.Stderr, "Using fieldname: %s\n", fieldName)

	queryStringRaw := fmt.Sprintf(`SELECT Id, ParentId, Parent.Profile.Name, Field
		FROM FieldPermissions
		WHERE Parent.Type = 'Profile'
		AND Field='%s'
		ORDER BY Parent.Profile.Name`, fieldName)
	queryString := tidyForQueryParam(queryStringRaw)
	targetUrl := fmt.Sprintf("%s/services/data/v%s/query?q=%s", sfInfo.instanceUrl, sfInfo.version, queryString)
	fmt.Fprintf(os.Stderr, "Using url: %s\n", targetUrl)
	fetchResp, fetchRespErr := getWrapper.AuthedGet(targetUrl, sfInfo.accessToken)
	if fetchRespErr != nil {
		panic("Error executing GET Request: ")
	}

	os.Stderr.WriteString(sfInfo.instanceUrl)
	os.Stderr.WriteString(fetchResp.Status)

	fetchRespStr, fetchRespStrErr := io.ReadAll(fetchResp.Body)
	//fmt.Println(string(fetchRespStr))
	if fetchRespStrErr != nil {
		panic("Error reading GET Request")
	}
	var currentPermsResult sfQueryResult[sfFieldPermissons]
	currentPermsParseErr := json.Unmarshal(fetchRespStr, &currentPermsResult)
	if currentPermsParseErr != nil {
		return fmt.Errorf("error parsing GET Request. Error: %w", currentPermsParseErr)
	}

	if !currentPermsResult.Done {
		panic("Incomplete result set. Extra result parsing not implemented.")
	}
	if len(currentPermsResult.Records) == 0 {
		os.Stderr.WriteString("No records")
		return nil
	}

	currentPermsList := currentPermsResult.Records

	outputPerms := []string{}
	for _, p := range currentPermsList {
		pKey := p.ParentId
		permString := "N"
		outputPerms = append(outputPerms, fmt.Sprintf("%s:%s", pKey, permString))
	}

	outputPermsString := strings.Join(outputPerms, ";")

	fmt.Println(outputPermsString)

	return nil
}
