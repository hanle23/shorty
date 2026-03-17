package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/hanle23/shorty/internal/config"
	"github.com/hanle23/shorty/internal/types"
)

func setupTestEnv(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	config.ResetForTesting()
	return tmpDir
}

func pipeStdin(t *testing.T, input string) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	w.WriteString(input)
	w.Close()
	oldStdin := os.Stdin
	os.Stdin = r
	t.Cleanup(func() { os.Stdin = oldStdin })
}

func configDir(tmpDir string) string {
	return filepath.Join(tmpDir, ".config", "shorty")
}

func TestInitFlow_FirstTime(t *testing.T) {
	tmpDir := setupTestEnv(t)
	dir := configDir(tmpDir)

	err := config.InitFlow()
	if err != nil {
		t.Fatalf("InitFlow failed: %v", err)
	}

	configPath := filepath.Join(dir, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config.yaml was not created")
	}

	runnablePath := filepath.Join(dir, "runnables.yaml")
	if _, err := os.Stat(runnablePath); os.IsNotExist(err) {
		t.Fatal("runnables.yaml was not created")
	}

	// Verify config content points to the correct runnables path
	configData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config.yaml: %v", err)
	}
	var cfg types.ConfigFile
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		t.Fatalf("failed to parse config.yaml: %v", err)
	}
	if cfg.RunnablePath != runnablePath {
		t.Errorf("config runnable_path = %q, want %q", cfg.RunnablePath, runnablePath)
	}

	// Verify runnables content is empty
	runnableData, err := os.ReadFile(runnablePath)
	if err != nil {
		t.Fatalf("failed to read runnables.yaml: %v", err)
	}
	var runnable types.RunnableFile
	if err := yaml.Unmarshal(runnableData, &runnable); err != nil {
		t.Fatalf("failed to parse runnables.yaml: %v", err)
	}
	if len(runnable.Shortcuts) != 0 {
		t.Errorf("expected 0 shortcuts, got %d", len(runnable.Shortcuts))
	}
	if len(runnable.Scripts) != 0 {
		t.Errorf("expected 0 scripts, got %d", len(runnable.Scripts))
	}
}

func TestInitFlow_ResetConfirm(t *testing.T) {
	tmpDir := setupTestEnv(t)
	dir := configDir(tmpDir)

	// Run init once to create files
	if err := config.InitFlow(); err != nil {
		t.Fatalf("first InitFlow failed: %v", err)
	}

	// Write some data into runnables to simulate user modifications
	runnablePath := filepath.Join(dir, "runnables.yaml")
	modified := types.RunnableFile{
		Shortcuts: map[string]types.Shortcut{
			"test": {Shortcut_name: "test", Package_name: "echo", Args: []string{"hello"}},
		},
		Scripts: make(map[string]types.Script),
	}
	modifiedYaml, _ := yaml.Marshal(modified)
	if err := os.WriteFile(runnablePath, modifiedYaml, 0644); err != nil {
		t.Fatalf("failed to write modified runnables: %v", err)
	}

	// Reset cache so InitFlow re-reads from disk
	config.ResetForTesting()

	// Pipe "y" to confirm reset
	pipeStdin(t, "y\n")

	if err := config.InitFlow(); err != nil {
		t.Fatalf("reset InitFlow failed: %v", err)
	}

	// Verify runnables were reset to empty
	runnableData, err := os.ReadFile(runnablePath)
	if err != nil {
		t.Fatalf("failed to read runnables.yaml: %v", err)
	}
	var runnable types.RunnableFile
	if err := yaml.Unmarshal(runnableData, &runnable); err != nil {
		t.Fatalf("failed to parse runnables.yaml: %v", err)
	}
	if len(runnable.Shortcuts) != 0 {
		t.Errorf("expected 0 shortcuts after reset, got %d", len(runnable.Shortcuts))
	}
}

func TestInitFlow_ResetDecline(t *testing.T) {
	tmpDir := setupTestEnv(t)
	dir := configDir(tmpDir)

	// Run init once to create files
	if err := config.InitFlow(); err != nil {
		t.Fatalf("first InitFlow failed: %v", err)
	}

	// Write some data into runnables to simulate user modifications
	runnablePath := filepath.Join(dir, "runnables.yaml")
	modified := types.RunnableFile{
		Shortcuts: map[string]types.Shortcut{
			"test": {Shortcut_name: "test", Package_name: "echo", Args: []string{"hello"}},
		},
		Scripts: make(map[string]types.Script),
	}
	modifiedYaml, _ := yaml.Marshal(modified)
	if err := os.WriteFile(runnablePath, modifiedYaml, 0644); err != nil {
		t.Fatalf("failed to write modified runnables: %v", err)
	}

	// Reset cache so InitFlow re-reads from disk
	config.ResetForTesting()

	// Pipe "n" to decline reset
	pipeStdin(t, "n\n")

	if err := config.InitFlow(); err != nil {
		t.Fatalf("decline InitFlow failed: %v", err)
	}

	// Verify runnables were NOT reset — the shortcut should still be there
	runnableData, err := os.ReadFile(runnablePath)
	if err != nil {
		t.Fatalf("failed to read runnables.yaml: %v", err)
	}
	var runnable types.RunnableFile
	if err := yaml.Unmarshal(runnableData, &runnable); err != nil {
		t.Fatalf("failed to parse runnables.yaml: %v", err)
	}
	if len(runnable.Shortcuts) != 1 {
		t.Errorf("expected 1 shortcut (unchanged), got %d", len(runnable.Shortcuts))
	}
	if _, exists := runnable.Shortcuts["test"]; !exists {
		t.Error("expected 'test' shortcut to still exist after declining reset")
	}
}

func TestSetRunnableDir_StoresFilePath(t *testing.T) {
	tmpDir := setupTestEnv(t)
	dir := configDir(tmpDir)

	// Initialize config first so SetRunnableDir can load it
	if err := config.InitFlow(); err != nil {
		t.Fatalf("InitFlow failed: %v", err)
	}
	config.ResetForTesting()

	customDir := filepath.Join(tmpDir, "custom", "runnables")
	if err := os.MkdirAll(customDir, 0755); err != nil {
		t.Fatalf("failed to create custom dir: %v", err)
	}

	if err := config.SetRunnableDir(customDir); err != nil {
		t.Fatalf("SetRunnableDir failed: %v", err)
	}

	// Re-read config from disk and verify RunnablePath is a file path, not a directory
	config.ResetForTesting()
	configData, err := os.ReadFile(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatalf("failed to read config.yaml: %v", err)
	}
	var cfg types.ConfigFile
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		t.Fatalf("failed to parse config.yaml: %v", err)
	}
	expected := filepath.Join(customDir, "runnables.yaml")
	if cfg.RunnablePath != expected {
		t.Errorf("RunnablePath = %q, want %q", cfg.RunnablePath, expected)
	}
}

func TestLoadConfigAndRunnable_AfterInit(t *testing.T) {
	setupTestEnv(t)

	if err := config.InitFlow(); err != nil {
		t.Fatalf("InitFlow failed: %v", err)
	}

	// Clear cache to force reading from disk
	config.ResetForTesting()

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed after init: %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig returned nil after init")
	}
	if cfg.RunnablePath == "" {
		t.Error("LoadConfig returned empty RunnablePath after init")
	}

	// Clear cache again to test GetRunnable independently
	config.ResetForTesting()

	runnable, err := config.GetRunnable()
	if err != nil {
		t.Fatalf("GetRunnable failed after init: %v", err)
	}
	if runnable == nil {
		t.Fatal("GetRunnable returned nil after init")
	}
	if runnable.Shortcuts == nil {
		t.Error("GetRunnable returned nil Shortcuts map")
	}
	if runnable.Scripts == nil {
		t.Error("GetRunnable returned nil Scripts map")
	}
}

func TestSaveYamlFile_OverwritesExistingContent(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.yaml")

	first := []byte("key: first_value\nother: data\n")
	if err := config.SaveYamlFile(filePath, first); err != nil {
		t.Fatalf("first SaveYamlFile failed: %v", err)
	}

	second := []byte("key: second\n")
	if err := config.SaveYamlFile(filePath, second); err != nil {
		t.Fatalf("second SaveYamlFile failed: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(got) != string(second) {
		t.Errorf("file content = %q, want %q (old content was not fully overwritten)", string(got), string(second))
	}
}
