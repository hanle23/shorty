/*
Copyright © 2025 Han Le <hanle.cs23@gmail.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/hanle23/shorty/config"
	"github.com/hanle23/shorty/internal/helper"

	"github.com/spf13/cobra"
	"os"
	"strings"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config or reset config to original state.",
	Run: func(cmd *cobra.Command, args []string) {
		currConfigDir := config.Dir()
		isExist := helper.IsExist(currConfigDir)
		r := bufio.NewReader(os.Stdin)
		if isExist {

			fmt.Printf("Found an existing configuration file (%s), do you want to override this? (y/n)? ", currConfigDir)
			ans, _ := r.ReadString('\n')
			ans = strings.TrimSpace(ans)
			if ans == "n" {
				return
			}
		}
		defaultPath := config.DefaultPath()
		fmt.Printf("Do you want to override the default path? (%s) (y/n): ", defaultPath)
		ans, _ := r.ReadString('\n')
		ans = strings.TrimSpace(ans)
		if ans == "n" {
			fmt.Println("Initiating config to default path...")
			err := config.SetDefaultConfigDir()
			cobra.CheckErr(err)
			return
		}

		fmt.Print("Please type the full path for the new config file: ")
		ans, _ = r.ReadString('\n')
		ans = strings.TrimSpace(ans)
		fmt.Println("Initiating config to overrided path...")
		err := config.SetOverrideConfigDir(ans)
		cobra.CheckErr(err)

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
