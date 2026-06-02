package mcpkit

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// claudeName is the name used when registering with the `claude` CLI.
// It strips the conventional "mcp-" prefix from the app name:
//   "mcp-atlassian" -> "atlassian".
func (a App) claudeName() string {
	return strings.TrimPrefix(a.Name, "mcp-")
}

// binaryPath returns the absolute, symlink-resolved path to the running binary.
func binaryPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	return filepath.Abs(exe)
}

func (a App) registerInClaude() error {
	exe, err := binaryPath()
	if err != nil {
		return fmt.Errorf("resolving binary path: %w", err)
	}

	name := a.claudeName()

	// Remove any existing registration first so `add` is idempotent.
	// Ignore error — may simply not be registered yet.
	_ = exec.Command("claude", "mcp", "remove", name).Run()

	cmd := exec.Command("claude", "mcp", "add", name, "--", exe)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("claude mcp add %s: %w", name, err)
	}
	fmt.Fprintf(os.Stderr, "Registered %q -> %s\n", name, exe)
	return nil
}

func (a App) unregisterFromClaude() error {
	name := a.claudeName()
	cmd := exec.Command("claude", "mcp", "remove", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("claude mcp remove %s: %w", name, err)
	}
	return nil
}

const opencodeConfig = "opencode.json"

func (a App) registerInOpencode() error {
	exe, err := binaryPath()
	if err != nil {
		return fmt.Errorf("resolving binary path: %w", err)
	}

	cfg, err := readOpencodeConfig()
	if err != nil {
		return err
	}

	mcp, _ := cfg["mcp"].(map[string]any)
	if mcp == nil {
		mcp = map[string]any{}
	}
	mcp[a.claudeName()] = map[string]any{
		"type":    "local",
		"command": []string{exe},
	}
	cfg["mcp"] = mcp

	if err := writeOpencodeConfig(cfg); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Registered %q -> %s in %s\n", a.claudeName(), exe, opencodeConfig)
	return nil
}

func (a App) unregisterFromOpencode() error {
	cfg, err := readOpencodeConfig()
	if err != nil {
		return err
	}

	mcp, _ := cfg["mcp"].(map[string]any)
	if mcp == nil {
		return fmt.Errorf("%q not found in %s", a.claudeName(), opencodeConfig)
	}
	if _, ok := mcp[a.claudeName()]; !ok {
		return fmt.Errorf("%q not found in %s", a.claudeName(), opencodeConfig)
	}
	delete(mcp, a.claudeName())
	if len(mcp) == 0 {
		delete(cfg, "mcp")
	} else {
		cfg["mcp"] = mcp
	}

	if err := writeOpencodeConfig(cfg); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Unregistered %q from %s\n", a.claudeName(), opencodeConfig)
	return nil
}

func readOpencodeConfig() (map[string]any, error) {
	data, err := os.ReadFile(opencodeConfig)
	if os.IsNotExist(err) {
		return map[string]any{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", opencodeConfig, err)
	}
	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", opencodeConfig, err)
	}
	return cfg, nil
}

func writeOpencodeConfig(cfg map[string]any) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	if err := os.WriteFile(opencodeConfig, append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", opencodeConfig, err)
	}
	return nil
}
