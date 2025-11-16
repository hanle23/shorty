/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/hanle23/shorty/internal/config"
	"github.com/hanle23/shorty/internal/strings"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Listing current shortcut and script.",
	Long: `List out every shortcut and script.

	This command will print out a non-ordered list of items from both shortcut and script,
	With every information from the user-set name, original name, description and additional argument
	per user selection.`,
	PersistentPreRun: func(cmd *cobra.Command, arg []string) {
		err := config.LoadRunnable()
		if err != nil {
			fmt.Println("Shorty list have error while loading shortcut: ")
			fmt.Println(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		shortcut, err := config.GetRunnable()
		if err != nil {
			fmt.Println("Shorty list have error while loading shortcut: ")
			fmt.Println(err)
			return
		}
		strings.PrintRunnable(*shortcut, 20)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
