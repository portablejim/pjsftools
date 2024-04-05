package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func GetPerms(args []string, runner CommandRunner, getWrapper AuthHttpGetter) error {

	fs := flag.NewFlagSet("getperms", flag.ExitOnError)
	useNames := fs.Bool("names", false, "use Profile Names instead of Ids")
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

	currentPermsList, getErr := getFieldPermsForField(fieldName, sfInfo, getWrapper)
	if getErr != nil {
		return fmt.Errorf("error getting perms: %w", getErr)
	}

	outputPerms := []string{}
	for _, p := range currentPermsList {
		if p.Field == fieldName {
			outputPerms = append(outputPerms, generatePermString(p, *useNames))
		}
	}

	outputPermsString := strings.Join(outputPerms, ";")

	fmt.Println(outputPermsString)

	return nil
}

func generatePermString(p sfFieldPermissons, useNames bool) string {
	pKey := p.ParentId
	if useNames {
		pKey = url.QueryEscape(p.Parent.Profile.Name)
	}

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
