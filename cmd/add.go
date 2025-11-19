/*
Copyright © 2025 Han Le <hanle.cs23@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/hanle23/shorty/internal/config"
	"github.com/hanle23/shorty/internal/io"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, arg []string) {
		err := config.LoadRunnable()
		if err != nil {
			fmt.Println("Shorty list have error while loading shortcut: ")
			fmt.Println(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		mode := io.ModeSelectorPrompt()
		if mode == -1 {
			fmt.Println("Exiting...")
			return
		}
		if mode == 0 {
			newShortcut := io.AddNewShortcutPrompt(args)
			fmt.Println(newShortcut)
		} else {
			newScript := io.AddNewScriptPrompt(args)
			fmt.Println(newScript)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
