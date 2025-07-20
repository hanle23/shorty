package cmd

import (
	"fmt"
	"github.com/hanle23/shorty/internal/config"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{

	Use:   "doctor",
	Short: "Diagnose shorty settings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🩺 Running diagnostics...\n")

		dir := config.Dir()
		if dir == "" {
			fmt.Println("❌ Failed to read current config dir")
		}
		fmt.Printf("✔ Current config directory: %s\n", dir)
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
