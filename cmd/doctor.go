package cmd

import (
	"fmt"
	"os/exec"

	"github.com/hanle23/shorty/internal/config"
	"github.com/hanle23/shorty/internal/fs"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose shorty settings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running diagnostics...")
		healthy := true

		// 1. Config directory
		configDir, err := config.GetDefaultFolderPath()
		if err != nil {
			fmt.Printf("  ✗ Config directory: failed to resolve path: %v\n", err)
			healthy = false
		} else if !fs.IsExist(configDir) {
			fmt.Printf("  ✗ Config directory: %s does not exist\n", configDir)
			healthy = false
		} else {
			fmt.Printf("  ✓ Config directory: %s\n", configDir)
		}

		// 2. Config file
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("  ✗ Config file: %v\n", err)
			healthy = false
		} else {
			fmt.Println("  ✓ Config file: loaded successfully")
		}

		// 3. Runnable path set in config
		runnablePath := ""
		if cfg != nil {
			runnablePath = cfg.RunnablePath
		}
		if runnablePath == "" {
			fmt.Println("  ✗ Runnable path: not set in config")
			healthy = false
		} else {
			fmt.Printf("  ✓ Runnable path: %s\n", runnablePath)
		}

		// 4. Shell availability
		shPath, err := exec.LookPath("sh")
		if err != nil {
			fmt.Println("  ✗ Shell: 'sh' not found in PATH — scripts will not be able to run")
			healthy = false
		} else {
			fmt.Printf("  ✓ Shell: %s\n", shPath)
		}

		// 5. Runnable file
		if runnablePath != "" {
			if !fs.IsExist(runnablePath) {
				fmt.Printf("  ✗ Runnable file: %s does not exist\n", runnablePath)
				healthy = false
			} else {
				runnable, err := config.GetRunnable()
				if err != nil {
					fmt.Printf("  ✗ Runnable file: %v\n", err)
					healthy = false
				} else {
					fmt.Printf("  ✓ Runnable file: %d shortcut(s), %d script(s)\n",
						len(runnable.Shortcuts), len(runnable.Scripts))
				}
			}
		}

		fmt.Println()
		if healthy {
			fmt.Println("All checks passed.")
		} else {
			fmt.Println("Issues found. Run 'shorty init' to set up or reset your configuration.")
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
