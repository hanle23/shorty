/*
Copyright © 2025 Han Le <hanle.cs23@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config or reset config to original state.",
	Run: func(cmd *cobra.Command, args []string) {
		// err := config.InitFlow(true)
		// cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
