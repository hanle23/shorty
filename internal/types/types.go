package types

// All shortcut name needs to be unique
// Package Name is the name of the program that will be replaced by shortcut name
// Args is all the arguments that the user wants to include after package name
// Description is the string description of the shortcut during list command
type Shortcut struct {
	Shortcut_name string   `yaml:"shortcut_name"`
	Package_name  string   `yaml:"package_name"`
	Args          []string `yaml:"args"`
	Description   string   `yaml:"description,omitempty"`
}

// All package name needs to be unique
// Script is the actual script that will be triggered, including all the args
// Description is the string description of the shortcut during list command
type Script struct {
	Package_name string `yaml:"package_name"`
	Script       string `yaml:"script"`
	Description  string `yaml:"description,omitempty"`
}

type RunnableFile struct {
	Shortcuts map[string]Shortcut `yaml:"shortcuts"`
	Scripts   map[string]Script   `yaml:"scripts"`
}

type ConfigFile struct {
	RunnablePath string `yaml:"runnable_path"`
}
