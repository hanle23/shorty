/*
Copyright © 2025 Han Le <hanle.cs23@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/hanle23/shorty/internal/config"
	"github.com/hanle23/shorty/internal/fs"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config or reset config to original state.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		currShortcutDir := Dir()
		isExist := fs.IsExist(currShortcutDir)

	},
	Run: func(cmd *cobra.Command, args []string) {
		// err := config.InitFlow(true)
		// cobra.CheckErr(err)
		bytes, err := config.GetEmptyConfigYAML()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(bytes))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
