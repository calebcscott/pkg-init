package cmd

import (
	"github.com/calebcscott/pkg-init/pkg/config"
	"github.com/spf13/cobra"
)


var addCmd = &cobra.Command{
    Use: "add",
    Short: "Add a project template",
    Long: "Adds a project template so that it can be used shorthand.",
    Args: cobra.MaximumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        //TODO add logic to add template 

        // get viper configurations(?)
        _ = config.NewConfig()


        // parse/validate provided template

        // add template to config with viper
    },
}


func init() {
    rootCmd.AddCommand(addCmd)

    // TODO add flags
}
