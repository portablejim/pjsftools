package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func SetPerms(args []string, runner CommandRunner, getWrapper AuthHttpGetter, patchWrapper AuthHttpPatcher) error {

	fs := flag.NewFlagSet("setperms", flag.ExitOnError)
	useNames := fs.Bool("names", false, "use Profile Names instead of Ids")
	orgName := fs.String("org", "", "sf org to use")
	fs.Parse(args)

	otherArgs := fs.Args()
	if len(otherArgs) != 2 {
		return fmt.Errorf("expected exactly 2 args (field and perms). Got %v", otherArgs)
	}

	fmt.Println(*useNames)

	if len(*orgName) == 0 {
		os.Stderr.WriteString("setperms using default org\n")
	} else {
		os.Stderr.WriteString("setperms using org " + *orgName + "\n")
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

	// Parse field perms.
	fieldPerms := fs.Arg(1)
	permParts := strings.Split(fieldPerms, ";")
	permMap := map[string]string{}

	for i := 0; i < len(permParts); i++ {
		permPartFields := strings.Split(permParts[i], ":")
		if len(permPartFields) == 2 && len(permPartFields[0]) > 0 {
			// Has the correct formatting
			permPartFieldName := permPartFields[0]
			if *useNames {
				// Need to decode.
				tempFieldName, unescapeErr := url.QueryUnescape(permPartFieldName)
				if unescapeErr != nil {
					// Error parsing, skip
					continue
				}
				permPartFieldName = tempFieldName
			}
			permMap[permPartFieldName] = permPartFields[1]
		}

	}

	// Build changed perms.
	newPermList := []sfUpdateFieldPermissions{}
	changedPerms := []string{}
	for i := 0; i < len(currentPermsList); i++ {
		currentPerm := currentPermsList[i]
		permKey := currentPerm.ParentId
		if *useNames {
			permKey = currentPerm.Parent.Profile.Name
		}
		newPermStr := permMap[permKey]

		if len(newPermStr) == 0 || len(currentPerm.Id) == 0 {
			continue
		}

		newPerm := sfUpdateFieldPermissions{Id: currentPerm.Id}

		if newPermStr == "RW" || newPermStr == "R" {
			newPerm.PermissionsRead = true
		} else {
			newPerm.PermissionsRead = false
		}
		if newPermStr == "RW" {
			newPerm.PermissionsEdit = true
		} else {
			newPerm.PermissionsEdit = false
		}

		newPermList = append(newPermList, newPerm)
		changedPerms = append(changedPerms, permKey)
	}

	for i := 0; i < len(changedPerms); i++ {
		delete(permMap, changedPerms[i])
	}

	permsUpdated, permsUpdatedErr := setFieldPermsForField(newPermList, sfInfo, patchWrapper)
	if permsUpdatedErr != nil {
		return fmt.Errorf("error setting perms: %w", permsUpdatedErr)
	}

	fmt.Fprintf(os.Stderr, "%d permissions updated\n", len(permsUpdated))

	fmt.Printf("Perms: %v\n", newPermList)

	return nil
}
