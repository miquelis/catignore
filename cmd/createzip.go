/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"log"
	"path/filepath"

	createzip "github.com/miquelis/catignore/createZip"
	"github.com/spf13/cobra"
)

var catIgnorePath, catIgnoreName, outputPath, nameFileZipeed string
var err error

// createzipCmd represents the createzip command
var createzipCmd = &cobra.Command{
	Use:   "createzip",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

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
}

func init() {
	rootCmd.AddCommand(createzipCmd)

	createzipCmd.Flags().StringVarP(
		&catIgnorePath,
		"path-catignore",
		"p",
		".",
		"Specify the path of the .catignore file or one of the supported files. If not specified, the directory being executed will be used.",
	)

	createzipCmd.Flags().StringVarP(
		&catIgnoreName,
		"name-catignore",
		"c",
		".catignore",
		`Specify the name of the .ignore file. If it does not exist, the ".catignore" file will be used by default.
		".catignore" and ".gcloudignore" files are currently supported.`,
	)

	createzipCmd.Flags().StringVarP(
		&outputPath,
		"output",
		"o",
		"/tmp",
		`Specify the path where the zip file will be saved. If not specified, it will be saved in the directory "/tmp"`,
	)

	createzipCmd.Flags().StringVarP(
		&nameFileZipeed,
		"name-zip",
		"n",
		"functions",
		`Specify the name of the zip file. If not specified, the name "functions" will be added. No need to put the .zip extension`,
	)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createzipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createzipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
