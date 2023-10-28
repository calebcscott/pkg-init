/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/calebcscott/pkg-init/pkg/config"
	"github.com/calebcscott/pkg-init/pkg/template"
	"github.com/spf13/cobra"
)


// getCmd represents the get command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project template",
	Long: `Initialize a directory based on a template.

The default project template is default.yaml, but 
can be a separate template can be specified. The
default directory to install template to is the
current directory, but can be provided.`,
    Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

        templateFile, _ := cmd.Flags().GetString("template")
        lang, _ := cmd.Flags().GetString("lang")

        config := config.NewConfig()

        outputDir, _ := os.Getwd()
        if len(args) > 0 {
            outputDir = args[0]
        }

        fmt.Println("Using flags/options:", templateFile, lang, outputDir)

        template, err := template.NewTemplate(templateFile, lang, outputDir, &config)
        if err != nil {
            fmt.Println("Error", err)
            return
        }

        template.Init(&config)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    initCmd.Flags().StringP("template", "t", "default", "Template file to use instead of default")
    initCmd.Flags().StringP("lang", "l", "", "If there is a specific language default project template to use")

    initCmd.MarkFlagsMutuallyExclusive("template", "lang")
}
