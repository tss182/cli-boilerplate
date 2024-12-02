package lib

import (
	"fmt"
	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"regexp"
	"strings"
)

func PathExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Printf("Error checking folder existence: %v\n", err)
		return false
	}
	if info.Name() != "" {
		return true
	}
	return false
}

func RenameFile(s string) string {
	reg, _ := regexp.Compile("[^a-zA-Z]")
	var str = reg.ReplaceAllString(s, "-")
	reg, _ = regexp.Compile("-+")
	str = reg.ReplaceAllString(str, "-")
	str = strings.Trim(str, "-")
	return strings.ToLower(str)
}

func RenamePackage(s string) string {
	t := cases.Title(language.English)
	str := t.String(strings.TrimSpace(s))
	reg, _ := regexp.Compile("[^a-zA-Z]")
	str = reg.ReplaceAllString(s, "")
	return str
}

func GetModuleName() string {
	goModBytes, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Println("error reading go.mod: ", err.Error())
	}

	modName := modfile.ModulePath(goModBytes)

	return modName
}
