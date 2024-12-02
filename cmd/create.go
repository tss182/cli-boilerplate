package cmd

import (
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

		fmt.Println("module name", goModName)
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

		// Create the folder
		err := os.MkdirAll(domainFolder, os.ModePerm)
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
		_, err = handlerFile.WriteString("package " + lib.RenamePackage(domain))
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
		_, err = repoFile.WriteString("package repository")
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
		_, err = entityFile.WriteString("package entity")
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
		_, err = featFile.WriteString("package feature")
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
