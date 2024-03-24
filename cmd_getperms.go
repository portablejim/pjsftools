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
		os.Stderr.WriteString("getperms using default org\n")
	} else {
		os.Stderr.WriteString("getperms using org " + *orgName + "\n")
	}
	sfInfo, tokenError := getAccessToken(*orgName, runner)
	if tokenError != nil {
		return fmt.Errorf("unable to get access token. Error: %w", tokenError)
	}
	if len(sfInfo.accessToken) < 1 || len(sfInfo.instanceUrl) < 1 || len(sfInfo.version) < 1 {
		return fmt.Errorf("something is wrong. Unable to get valid info")
	}

	fieldName := fs.Arg(0)

	queryStringRaw := fmt.Sprintf(`SELECT Id, ParentId, Parent.Profile.Name, Field, PermissionsEdit, PermissionsRead
		FROM FieldPermissions
		WHERE Parent.Type = 'Profile'
		AND Field='%s'
		ORDER BY Parent.Profile.Name`, fieldName)
	queryString := tidyForQueryParam(queryStringRaw)
	targetUrl := fmt.Sprintf("%s/services/data/v%s/query?q=%s", sfInfo.instanceUrl, sfInfo.version, queryString)
	fetchResp, fetchRespErr := getWrapper.AuthedGet(targetUrl, sfInfo.accessToken)
	if fetchRespErr != nil {
		panic("Error executing GET Request")
	}

	fetchRespStr, fetchRespStrErr := io.ReadAll(fetchResp.Body)
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
		outputPerms = append(outputPerms, generatePermString(p))
	}

	outputPermsString := strings.Join(outputPerms, ";")

	fmt.Println(outputPermsString)

	return nil
}

func generatePermString(p sfFieldPermissons) string {
	pKey := p.ParentId

	permString := ""
	if p.PermissionsRead {
		permString += "R"
	}
	if p.PermissionsEdit {
		permString += "W"
	}
	if len(permString) == 0 {
		permString = "N"
	}

	return fmt.Sprintf("%s:%s", pKey, permString)
}
