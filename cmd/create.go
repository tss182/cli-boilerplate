package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tss182/cli-boilerplate/cmd/lib"
	"os"
	"path/filepath"
	"strings"
)

var (
	domain string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new folder and file with custom content",
	Long: `The create command allows you to generate a new folder and file.
You can specify the folder name, file name, and the content to write to the file.`,
	Run: func(cmd *cobra.Command, args []string) {
		//folder required
		folderRequired := []string{
			"./domain",
			"./delivery",
			"./delivery/container",
			"./delivery/container/container.go",
			"./delivery/server",
			"./delivery/server/handler.go",
			"./delivery/server/router_v1.go",
			"./main.go",
			"./go.mod",
		}

		goModName := lib.GetModuleName()

		//check folder domain
		for _, v := range folderRequired {
			if !lib.PathExists(v) {
				fmt.Printf("Path %s does not exist!\n", v)
				return
			}
		}

		// Ensure folder name is provided
		if strings.TrimSpace(domain) == "" {
			fmt.Println("Error: domain name is required")
			return
		}

		domainFolder := filepath.Join("./domain", lib.RenameFile(domain))
		if lib.PathExists(domainFolder) {
			fmt.Printf("Domain %s already exists!\n", domainFolder)
			return
		}

		//Template value
		var templateValue lib.ValueTemplate
		templateValue.GoModName = goModName
		templateValue.DomainPackage = lib.PackageName(domain, true, "")
		templateValue.DomainPackageLocal = lib.PackageName(domain, false, "")
		templateValue.Domain = lib.PackageName(domain, true, " ")
		templateValue.Folder = strings.Replace(domainFolder, "\\", "/", -1)

		//container.go
		targetFileContainer := "./delivery/container/container.go"
		file, err := os.OpenFile(targetFileContainer, os.O_RDWR, 0755)
		if err != nil {
			fmt.Printf("failed to open file: %v", err)
			return
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		var content []string
		scanner := bufio.NewScanner(file)

		//import config
		inImport := false
		newImportFeat := fmt.Sprintf("\t%sFeature \"%s/%s/feature\"", templateValue.DomainPackageLocal, templateValue.GoModName, templateValue.Folder)
		newImportRepo := fmt.Sprintf("\t%sRepository \"%s/%s/repository\"", templateValue.DomainPackageLocal, templateValue.GoModName, templateValue.Folder)

		//struct config
		inStruct := false
		newStruct := fmt.Sprintf("\t%sFeature %sFeature.%sFeatureInterface", templateValue.DomainPackage, templateValue.DomainPackageLocal, templateValue.DomainPackage)

		//container return config
		inContainerRetrun := false
		newVariableRepo := fmt.Sprintf("\t%sRepo := %sRepository.New(dbMySQL,&cfg)", templateValue.DomainPackageLocal, templateValue.DomainPackageLocal)
		newVariableFeat := fmt.Sprintf("\t%sFeat := %sFeature.New(%sRepo)", templateValue.DomainPackageLocal, templateValue.DomainPackageLocal, templateValue.DomainPackageLocal)
		newValueStruct := fmt.Sprintf("\t\t%sFeature: %sFeat,", templateValue.DomainPackage, templateValue.DomainPackageLocal)
		//ApprovalCorrectionFeature: approvalCorrectionFeat,
		for scanner.Scan() {
			line := scanner.Text()
			lineNoSpace := strings.TrimSpace(strings.Replace(line, " ", "", -1))
			if strings.HasPrefix(lineNoSpace, "import(") {
				inImport = true
			}

			if inImport && strings.HasPrefix(lineNoSpace, ")") {
				inImport = false
				content = append(content, newImportFeat)
				content = append(content, newImportRepo)
			}

			if strings.HasPrefix(lineNoSpace, "typeContainerstruct{") {
				inStruct = true
			}

			if inStruct && strings.HasPrefix(lineNoSpace, "}") {
				inStruct = false
				content = append(content, newStruct)
			}

			if strings.HasPrefix(lineNoSpace, "returnContainer{") {
				content = append(content, newVariableRepo)
				content = append(content, newVariableFeat)
				content = append(content, "\n\n")
				inContainerRetrun = true
			}

			if inContainerRetrun && strings.HasPrefix(lineNoSpace, "}") {
				inContainerRetrun = false
				content = append(content, newValueStruct)

			}

			content = append(content, line)
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("failed to read file: %v", err)
			return
		}

		// Write the modified content back to the file
		_ = file.Truncate(0)   // Clear the file
		_, _ = file.Seek(0, 0) // Move pointer to the start
		writer := bufio.NewWriter(file)
		for _, line := range content {
			_, _ = writer.WriteString(line + "\n")
		}
		_ = writer.Flush()

		//handler.go
		targetFileHandler := "./delivery/server/handler.go"
		fileHandler, err := os.OpenFile(targetFileHandler, os.O_RDWR, 0755)
		if err != nil {
			fmt.Printf("failed to open file: %v", err)
			return
		}
		defer func(fileHandler *os.File) {
			_ = fileHandler.Close()
		}(fileHandler)

		content = []string{}
		scanner = bufio.NewScanner(fileHandler)

		//import config
		inImport = false
		newImport := fmt.Sprintf("\t\"%s/%s\"", templateValue.GoModName, templateValue.Folder)

		//struct config
		inStruct = false
		newStruct = fmt.Sprintf("\t%s %s.%sHandlerInterface", templateValue.DomainPackage, templateValue.DomainPackageLocal, templateValue.DomainPackage)

		//value struct config
		inValueStruct := false
		newVariableStruct := fmt.Sprintf("\t\t%s : \t%s.New(cont.%sFeature),", templateValue.DomainPackage, templateValue.DomainPackageLocal, templateValue.DomainPackage)

		for scanner.Scan() {
			line := scanner.Text()
			lineNoSpace := strings.TrimSpace(strings.Replace(line, " ", "", -1))
			if strings.HasPrefix(lineNoSpace, "import(") {
				inImport = true
			}

			if inImport && strings.HasPrefix(lineNoSpace, ")") {
				inImport = false
				content = append(content, newImport)
			}

			if strings.HasPrefix(lineNoSpace, "typeHandlerstruct{") {
				inStruct = true
			}

			if inStruct && strings.HasPrefix(lineNoSpace, "}") {
				inStruct = false
				content = append(content, newStruct)
			}

			if strings.HasPrefix(lineNoSpace, "returnHandler{") {
				inValueStruct = true
			}

			if inValueStruct && strings.HasPrefix(lineNoSpace, "}") {
				inValueStruct = false
				content = append(content, newVariableStruct)
			}

			content = append(content, line)
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("failed to read file: %v", err)
			return
		}

		// Write the modified content back to the file
		_ = fileHandler.Truncate(0)   // Clear the file
		_, _ = fileHandler.Seek(0, 0) // Move pointer to the start
		writer = bufio.NewWriter(fileHandler)
		for _, line := range content {
			_, _ = writer.WriteString(line + "\n")
		}
		_ = writer.Flush()

		// Create the folder
		err = os.MkdirAll(domainFolder, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating folder: %v\n", err)
			return
		}

		//create entity
		entityFolder := filepath.Join(domainFolder, "entity")
		err = os.MkdirAll(entityFolder, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating folder: %v\n", err)
			return
		}

		//create feature
		featureFolder := filepath.Join(domainFolder, "feature")
		err = os.MkdirAll(featureFolder, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating folder: %v\n", err)
			return
		}

		//create repository
		repositoryFolder := filepath.Join(domainFolder, "repository")
		err = os.MkdirAll(repositoryFolder, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating folder: %v\n", err)
			return
		}

		//create handler file
		handler := filepath.Join(domainFolder, "handler.go")
		handlerFile, err := os.Create(handler)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}
		defer func(handlerFile *os.File) {
			err := handlerFile.Close()
			if err != nil {
				fmt.Printf("Error closing file: %v\n", err)
			}
		}(handlerFile)

		// Write content to the file
		handlerContent, err := lib.TemplateParse(lib.SampleHandler, templateValue)
		if err != nil {
			fmt.Printf("Error parsing template: %v\n", err)
			return
		}
		_, err = handlerFile.WriteString(handlerContent)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}

		//create repo file
		repo := filepath.Join(repositoryFolder, "repository.go")
		repoFile, err := os.Create(repo)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}
		defer func(repoFile *os.File) {
			err := repoFile.Close()
			if err != nil {
				fmt.Printf("Error closing file: %v\n", err)
			}
		}(repoFile)

		// Write content to the file
		repoContent, err := lib.TemplateParse(lib.SampleRepository, templateValue)
		if err != nil {
			fmt.Printf("Error parsing template: %v\n", err)
			return
		}
		_, err = repoFile.WriteString(repoContent)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}

		//create entity file
		entity := filepath.Join(entityFolder, "entity.go")
		entityFile, err := os.Create(entity)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}
		defer func(entityFile *os.File) {
			err := entityFile.Close()
			if err != nil {
				fmt.Printf("Error closing file: %v\n", err)
			}
		}(entityFile)

		// Write content to the file
		entityContent, err := lib.TemplateParse(lib.SampleEntity, templateValue)
		if err != nil {
			fmt.Printf("Error parsing template: %v\n", err)
			return
		}
		_, err = entityFile.WriteString(entityContent)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}

		//create feature file
		feat := filepath.Join(featureFolder, "feature.go")
		featFile, err := os.Create(feat)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}
		defer func(featFile *os.File) {
			err := featFile.Close()
			if err != nil {
				fmt.Printf("Error closing file: %v\n", err)
			}
		}(featFile)

		// Write content to the file
		featContent, err := lib.TemplateParse(lib.SampleFeature, templateValue)
		if err != nil {
			fmt.Printf("Error parsing template: %v\n", err)
			return
		}
		_, err = featFile.WriteString(featContent)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}

		fmt.Printf("Domain '%s' created successfully.\n", domain)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Add flags
	createCmd.Flags().StringVarP(&domain, "domain", "d", "", "Name of the domain to create")
}
