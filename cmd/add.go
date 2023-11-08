package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/calebcscott/pkg-init/pkg/config"
	"github.com/calebcscott/pkg-init/pkg/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


var addCmd = &cobra.Command{
    Use: "add",
    Short: "Add a project template",
    Long: "Adds a project template so that it can be used shorthand.",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        //TODO add logic to add template 
        templateName, _ := cmd.Flags().GetString("name")
        _, fileName := filepath.Split(args[0])
        if templateName == "" {
            templateName = strings.Split(fileName, ".")[0]
        }
        lang, _ := cmd.Flags().GetString("lang")

        fmt.Println("Adding template:", templateName)
        fmt.Println("for language type:", lang)

        // get viper configurations(?)
        config := config.NewConfig()

        // attempt to copy/cache template file
        dstPath := template.CacheTemplateFile(args[0], &config)

        // parse/validate provided template

        
        // add template to config with viper
        v := viper.GetString("templates."+templateName)
        if v != "" {
            if !template.GetChoiceB("Overrideing template: "+v+", Do you want to continue?[Y/n]") {
                return
            }

        }
        viper.Set("templates."+templateName, dstPath)

        err := viper.WriteConfig()
        if err != nil {
            fmt.Println("Error writing config")
        }
    },
}


func init() {
    rootCmd.AddCommand(addCmd)

    // TODO add name flag, this may have to be done after adjustment to yaml file reading
    //  to allow name of template to be different from name of file
    addCmd.Flags().StringP("name", "n", "", "Name for template") 
    addCmd.Flags().StringP("lang", "l", "", "Language/Project type for template") 
}
