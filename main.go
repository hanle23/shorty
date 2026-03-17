/*
Copyright © 2025 Han Le <hanle.cs23@gmail.com>
*/
package main

import (
	"github.com/hanle23/shorty/cmd"
)

var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
