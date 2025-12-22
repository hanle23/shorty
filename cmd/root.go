/*
Copyright © 2025 Han Le <hanle.cs23@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hanle23/shorty/internal/config"
	"github.com/hanle23/shorty/internal/context"
	"github.com/hanle23/shorty/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultLicense = "GPL-3.0"
)

var (
	cfgFile string
	debug   bool

	rootCmd = &cobra.Command{
		Use:   "shorty [name] [args...]",
		Short: "A CLI tool to manage and run shortcuts and scripts",
		Long: `Shorty is a CLI tool that helps you manage and execute shortcuts and scripts.

Usage:
  shorty <name> [args...]           Run a shortcut or script directly
  shorty shortcut <name> [args...]  Run a shortcut explicitly
  shorty script <name> [args...]    Run a script explicitly

When a name matches both a shortcut and a script, the shortcut takes precedence.
Use 'shorty script <name>' to explicitly run the script in case of collision.`,
		Args:                  cobra.ArbitraryArgs,
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			name := args[0]
			passedArgs := args[1:]

			runnables, err := config.GetRunnable()
			if err != nil {
				return fmt.Errorf("failed to load runnables: %w", err)
			}

			// Try shortcut first
			if shortcut, exists := runnables.Shortcuts[name]; exists {
				return executeShortcut(&shortcut, passedArgs)
			}

			// Fall back to script
			if script, exists := runnables.Scripts[name]; exists {
				return executeScript(&script, passedArgs)
			}

			return fmt.Errorf("no shortcut or script found with name: %s\nRun 'shorty --help' for usage", name)
		},
	}

	runCmd = &cobra.Command{
		Use:   "run <name> [args...]",
		Short: "Run a shortcut or script (shortcut takes precedence)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			passedArgs := args[1:]

			runnables, err := config.GetRunnable()
			if err != nil {
				return fmt.Errorf("failed to load runnables: %w", err)
			}

			// Try shortcut first
			if shortcut, exists := runnables.Shortcuts[name]; exists {
				return executeShortcut(&shortcut, passedArgs)
			}

			// Fall back to script
			if script, exists := runnables.Scripts[name]; exists {
				return executeScript(&script, passedArgs)
			}

			return fmt.Errorf("no shortcut or script found with name: %s", name)
		},
	}

	shortcutCmd = &cobra.Command{
		Use:     "shortcut <name> [args...]",
		Aliases: []string{"sc"},
		Short:   "Run a shortcut explicitly",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			passedArgs := args[1:]

			runnables, err := config.GetRunnable()
			if err != nil {
				return fmt.Errorf("failed to load runnables: %w", err)
			}

			shortcut, exists := runnables.Shortcuts[name]
			if !exists {
				return fmt.Errorf("shortcut not found: %s", name)
			}

			return executeShortcut(&shortcut, passedArgs)
		},
	}

	scriptCmd = &cobra.Command{
		Use:     "script <name> [args...]",
		Aliases: []string{"sr"},
		Short:   "Run a script explicitly",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			passedArgs := args[1:]

			runnables, err := config.GetRunnable()
			if err != nil {
				return fmt.Errorf("failed to load runnables: %w", err)
			}

			script, exists := runnables.Scripts[name]
			if !exists {
				return fmt.Errorf("script not found: %s", name)
			}

			return executeScript(&script, passedArgs)
		},
	}
)

func executeShortcut(shortcut *types.Shortcut, extraArgs []string) error {
	allArgs := append(shortcut.Args, extraArgs...)

	cmd := exec.Command(shortcut.Package_name, allArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if context.GetContext().Debug {
		fmt.Printf("Executing shortcut: %s %s\n", shortcut.Package_name, strings.Join(allArgs, " "))
	}

	return cmd.Run()
}

func executeScript(script *types.Script, extraArgs []string) error {
	scriptCmd := script.Script
	if len(extraArgs) > 0 {
		scriptCmd = scriptCmd + " " + strings.Join(extraArgs, " ")
	}

	cmd := exec.Command("sh", "-c", scriptCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if context.GetContext().Debug {
		fmt.Printf("Executing script: %s\n", scriptCmd)
	}

	return cmd.Run()
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initializeGlobalContextFromFlags)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "default", "config file (default is $HOME/.config/shorty/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Set debug mode")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.SetDefault("author", "Han Le <hanle.cs23@gmail.com>")
	viper.SetDefault("license", DefaultLicense)
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(shortcutCmd)
	rootCmd.AddCommand(scriptCmd)
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

// initializeGlobalContextFromFlags sets up the global context with parsed flag values
// This runs after flags are parsed but before command execution
func initializeGlobalContextFromFlags() {
	c := context.GetContext()
	c.SetContext(debug)
}
