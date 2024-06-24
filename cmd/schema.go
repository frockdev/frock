/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"frock2/podman"
	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
	"os"
)

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		createSchema()
		addSchemaFileToGitIgnore()
	},
}

// addSchemaFileToGitIgnore adds the schema file to the .gitignore file
// need to check first if file is already exists
// after that check if frock.schema.json is already in the file
// if not - create the file and add the line
// or just add the line if .gitignore is already created
func addSchemaFileToGitIgnore() {
	fileName := ".gitignore"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString("frock.schema.json\n")
		f.WriteString("frock.override.yaml\n")
	} else {
		//check if file already has it
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if scanner.Text() == "frock.schema.json" {
				return
			}
		}
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString("\nfrock.schema.json\n")
		f.WriteString("frock.override.yaml\n")
	}
}

func createSchema() {
	file := "frock.schema.json"
	var a = jsonschema.Reflect(&podman.Project{})
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	json, _ := a.MarshalJSON()
	f.WriteString(string(json))
}

func init() {
	rootCmd.AddCommand(schemaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// schemaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// schemaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
