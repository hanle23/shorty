/*
Copyright © 2025 Han Le <hanle.cs23@gmail.com>
*/
package cmd

import (
	"github.com/hanle23/shorty/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config or reset config to original state.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.InitFlow()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
