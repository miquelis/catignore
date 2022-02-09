/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	createzip "github.com/miquelis/catignore/createZip"
	"github.com/spf13/cobra"
)

var catIgnorePath, catIgnoreName, outputPath, nameFileZipeed string
var err error
var version bool = false
var VERSION string = "v0.0.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{

	Use:   "catignore",
	Short: color.YellowString("Catignore creates a .zip file checking which files (or patterns) to ignore."),
	Long: color.YellowString(`
Catignore creates a .zip file checking which files (or patterns) to ignore.
	
The files (or patterns) must be added in the .catignore file (or in the supported files list, example: ".gcloudignore")
	
For more information access the documentation: https://github.com/miquelis/catignore#readme`),

	Args: func(cmd *cobra.Command, args []string) error {

		if catIgnorePath == "." {
			catIgnorePath, err = filepath.Abs(catIgnorePath)

			if err != nil {
				return errors.New("invalid directory. Check and try again")
			}
			return nil
		}

		if catIgnoreName == " " {
			return errors.New("invalid .ignore file. Check and try again")
		}

		if outputPath != "/tmp" {
			path, err := filepath.Abs(outputPath)

			if err != nil {
				return errors.New("invalid directory. Check and try again")
			}

			outputPath = filepath.Join(path, "tmp")
			return nil
		}

		if nameFileZipeed == " " {
			return errors.New("invalid name file. Check and try again")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {

		if version {
			fmt.Println("catignore " + VERSION)
			os.Exit(0)
		}

		msg, err := createzip.CreateZipFile(
			filepath.Join(catIgnorePath, catIgnoreName),
			filepath.Join(outputPath, nameFileZipeed),
		)

		if err != nil {
			log.Fatal(err)
		}

		log.Println(msg)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// Added color in terminal cobra
	cobra.AddTemplateFunc("StyleHeading", color.New(color.FgGreen).SprintFunc())
	usageTemplate := rootCmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Aliases:`, `{{StyleHeading "Aliases:"}}`,
		`Available Commands:`, `{{StyleHeading "Available Commands:"}}`,
		`Global Flags:`, `{{StyleHeading "Global Flags:"}}`,
	).Replace(usageTemplate)

	re := regexp.MustCompile(`(?m)^Flags:\s*$`)
	usageTemplate = re.ReplaceAllLiteralString(usageTemplate, `{{StyleHeading "Flags:"}}`)

	rootCmd.SetUsageTemplate(usageTemplate)

	rootCmd.Flags().StringVarP(
		&catIgnorePath,
		"path-catignore",
		"p",
		".",
		color.CyanString("Specify the path of the .catignore file or one of the supported files. If not specified, the directory being executed will be used."),
	)

	rootCmd.Flags().StringVarP(
		&catIgnoreName,
		"name-catignore",
		"c",
		".catignore",
		color.CyanString(`Specify the name of the .ignore file. ".catignore" and ".gcloudignore" files are currently supported.`),
	)

	rootCmd.Flags().StringVarP(
		&outputPath,
		"output",
		"o",
		"/tmp",
		color.CyanString(`Specify the path where the zip file will be saved.`),
	)

	rootCmd.Flags().StringVarP(
		&nameFileZipeed,
		"name-zip",
		"n",
		"functions",
		color.CyanString(`Specify the name of the zip file. No need to put the .zip extension`),
	)

	rootCmd.Flags().BoolVarP(&version, "version", "v", false, color.CyanString("Print just the version number."))

}
