package lib

import (
	"bytes"
	"fmt"
	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type ValueTemplate struct {
	GoModName          string
	DomainPackage      string
	DomainPackageLocal string
	Domain             string
	Folder             string
}

func TemplateParse(txt string, data interface{}) (string, error) {
	tmpl, err := template.New("default").Parse(txt)
	if err != nil {
		return "", err
	}
	var parsedTemplate bytes.Buffer
	err = tmpl.Execute(&parsedTemplate, data)
	if err != nil {
		return "", err
	}

	return parsedTemplate.String(), nil
}

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

func PackageName(s string, upperFirstWord bool, separator string) string {
	reg, _ := regexp.Compile("[^a-zA-Z ]")
	s = reg.ReplaceAllString(s, "")
	t := cases.Title(language.English)
	arr := strings.Split(s, " ")
	var arrNew []string
	for i, v := range arr {
		if i == 0 && !upperFirstWord {
			continue
		}
		temp := t.String(strings.ToLower(strings.TrimSpace(v)))
		arrNew = append(arrNew, temp)
	}
	return strings.Join(arrNew, separator)
}

func GetModuleName() string {
	goModBytes, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Println("error reading go.mod: ", err.Error())
	}

	modName := modfile.ModulePath(goModBytes)

	return modName
}
