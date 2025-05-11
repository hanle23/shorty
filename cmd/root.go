/*
Copyright © 2025 Han Le <hanle.cs23@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/hanle23/shorty/config"
	"github.com/hanle23/shorty/internal/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "shorty [SHORTCUT] [ARGs...]",
		Short: "Run a shortcut or script",
		Args:  cobra.ArbitraryArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Grab value from config flag
			configDir, _ := cmd.Flags().GetString("config")
			if configDir != "default" {
				// Set overrided config if config is not empty
				err := config.SetOverrideConfigDir(configDir)
				cobra.CheckErr(err)
				fmt.Println("Successfully change config to: ", configDir)
			}
			// Check for current config file from normal flow
			currDir := config.Dir()
			isExist := helper.IsExist(currDir)
			// Run init flow to fix config issue if config does not exist
			if !isExist {
				fmt.Println("Config file was not successfully setup, init will be running now.")
				err := config.InitFlow()
				cobra.CheckErr(err)
			}

		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			shortcut := args[0]
			fmt.Println(shortcut)
			//TODO: Add configuration loader here, return error if config is not loaded
			fmt.Println("Root got called")
		},
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

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "default", "config file (default is $HOME/.config/shorty/config.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.SetDefault("author", "Han Le <hanle.cs23@gmail.com>")
	viper.SetDefault("license", "apache")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configPath := fmt.Sprintf("%s/.config/shorty", home)
		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}
	viper.AutomaticEnv()
}
