// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// PluginBinaryInfo contains information about a plugin binary
type PluginBinaryInfo struct {
	Name        string
	Path        string
	SHA256      string
	Version     string
	PluginType  string
	Environment map[string]string
}

// BuildPluginBinary builds a plugin binary from source code
func BuildPluginBinary(t *testing.T, sourcePath, outputPath string) error {
	t.Helper()

	// Ensure the output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
	}

	// Build the plugin binary
	cmd := exec.Command("go", "build", "-o", outputPath, sourcePath)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0") // Ensure static binary

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build plugin binary: %w\nOutput: %s", err, output)
	}

	t.Logf("Successfully built plugin binary: %s", outputPath)
	return nil
}

// DeployPluginBinary deploys a plugin binary to the plugin directory
func DeployPluginBinary(t *testing.T, binaryPath, pluginDir string) error {
	t.Helper()

	// Ensure the plugin directory exists
	if err := os.MkdirAll(pluginDir, 0o755); err != nil {
		return fmt.Errorf("failed to create plugin directory %s: %w", pluginDir, err)
	}

	// Copy the binary to the plugin directory
	binaryName := filepath.Base(binaryPath)
	targetPath := filepath.Join(pluginDir, binaryName)

	sourceFile, err := os.Open(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to open source binary %s: %w", binaryPath, err)
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create target binary %s: %w", targetPath, err)
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}

	// Make the binary executable
	if err := os.Chmod(targetPath, 0o755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	t.Logf("Successfully deployed plugin binary to: %s", targetPath)
	return nil
}

// ValidatePluginBinary validates a plugin binary
func ValidatePluginBinary(t *testing.T, binaryPath string) error {
	t.Helper()

	// Check if the file exists
	if _, err := os.Stat(binaryPath); err != nil {
		return fmt.Errorf("plugin binary not found at %s: %w", binaryPath, err)
	}

	// Check if the file is executable
	info, err := os.Stat(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to get binary info: %w", err)
	}

	if info.Mode()&0o111 == 0 {
		return fmt.Errorf("plugin binary %s is not executable", binaryPath)
	}

	// Try to run the binary with --help to see if it's a valid plugin
	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Note: Plugin binary %s may not support --help flag: %v", binaryPath, err)
	} else {
		t.Logf("Plugin binary %s validation output: %s", binaryPath, string(output))
	}

	t.Logf("Plugin binary validation passed for: %s", binaryPath)
	return nil
}

// GetPluginBinaryInfo extracts information about a plugin binary
func GetPluginBinaryInfo(t *testing.T, binaryPath string) (*PluginBinaryInfo, error) {
	t.Helper()

	// Calculate SHA256
	f, err := os.Open(binaryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open binary for SHA256 calculation: %w", err)
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return nil, fmt.Errorf("failed to calculate SHA256: %w", err)
	}

	sha256Sum := hex.EncodeToString(hasher.Sum(nil))
	binaryName := filepath.Base(binaryPath)

	info := &PluginBinaryInfo{
		Name:        strings.TrimSuffix(binaryName, filepath.Ext(binaryName)),
		Path:        binaryPath,
		SHA256:      sha256Sum,
		PluginType:  "unknown", // Will be determined by usage context
		Environment: make(map[string]string),
	}

	return info, nil
}

// ComparePluginVersions compares two plugin sessions for differences
func ComparePluginVersions(t *testing.T, v1, v2 *PluginSession) (*PluginDiff, error) {
	t.Helper()

	diff := &PluginDiff{
		PluginName: v1.PluginName,
		V1Info:     v1.GetPluginInfo(),
		V2Info:     v2.GetPluginInfo(),
		Changes:    make(map[string]interface{}),
	}

	// Compare binary paths (for external plugins)
	if !v1.IsBuiltin && !v2.IsBuiltin {
		if v1.BinaryPath != v2.BinaryPath {
			diff.Changes["binary_path"] = map[string]string{
				"old": v1.BinaryPath,
				"new": v2.BinaryPath,
			}
		}
	}

	// Compare configuration
	for key, v1Val := range v1.Config {
		if v2Val, exists := v2.Config[key]; !exists {
			diff.Changes["config_removed_"+key] = v1Val
		} else if v1Val != v2Val {
			diff.Changes["config_changed_"+key] = map[string]interface{}{
				"old": v1Val,
				"new": v2Val,
			}
		}
	}

	for key, v2Val := range v2.Config {
		if _, exists := v1.Config[key]; !exists {
			diff.Changes["config_added_"+key] = v2Val
		}
	}

	return diff, nil
}

// PluginDiff represents differences between plugin versions
type PluginDiff struct {
	PluginName string
	V1Info     map[string]interface{}
	V2Info     map[string]interface{}
	Changes    map[string]interface{}
}
