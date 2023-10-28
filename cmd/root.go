/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
    "fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)



// rootCmd represents the base command when called without any subcommands
var (
    cfgFile string

    rootCmd = &cobra.Command{
    	Use:   "pkg-init",
    	Short: "A tool to quickly create a project based off a template",
    	Long: ``,
    	// Uncomment the following line if your bare application
    	// has an action associated with it:
    	// Run: func(cmd *cobra.Command, args []string) { },
    }
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        configDir, err := os.UserConfigDir()
        cobra.CheckErr(err)

        var cmdDir = "/pkg-init"

        os.MkdirAll(configDir + cmdDir, 0755)

        viper.AddConfigPath(configDir + cmdDir)
        viper.AddConfigPath("/usr/share" + cmdDir + "/etc")
        viper.AddConfigPath(".")
        viper.SetConfigType("yaml")
        viper.SetConfigName("pkg-init")

    }

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err == nil {
        fmt.Println("User config file:", viper.ConfigFileUsed())
    }
}

func init() {
    cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pkg-init.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


