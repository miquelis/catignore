/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	createzip "github.com/miquelis/catignore/createZip"
	"github.com/spf13/cobra"
)

var catIgnorePath, catIgnoreName, outputPath, nameFileZipeed string
var err error

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "catignore",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	rootCmd.Flags().StringVarP(
		&catIgnorePath,
		"path-catignore",
		"p",
		".",
		"Specify the path of the .catignore file or one of the supported files. If not specified, the directory being executed will be used.",
	)

	rootCmd.Flags().StringVarP(
		&catIgnoreName,
		"name-catignore",
		"c",
		".catignore",
		`Specify the name of the .ignore file. ".catignore" and ".gcloudignore" files are currently supported.`,
	)

	rootCmd.Flags().StringVarP(
		&outputPath,
		"output",
		"o",
		"/tmp",
		`Specify the path where the zip file will be saved.`,
	)

	rootCmd.Flags().StringVarP(
		&nameFileZipeed,
		"name-zip",
		"n",
		"functions",
		`Specify the name of the zip file. No need to put the .zip extension`,
	)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.catignore.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
