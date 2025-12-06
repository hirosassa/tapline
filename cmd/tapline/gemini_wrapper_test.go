package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrapGemini_NoSessionManager(t *testing.T) {
	tmpDir := t.TempDir()
	invalidSessionDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(invalidSessionDir, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(invalidSessionDir, 0o750)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	mockGemini := filepath.Join(tmpDir, "gemini")
	mockScript := `#!/bin/bash
echo "Mock gemini output"
exit 0
`
	if err := os.WriteFile(mockGemini, []byte(mockScript), 0o755); err != nil {
		t.Fatal(err)
	}

	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir+":"+oldPath)
	defer os.Setenv("PATH", oldPath)

	wrapGemini([]string{"test"})
}

func TestWrapGemini_GeminiNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir)
	defer os.Setenv("PATH", oldPath)

	if os.Getenv("TEST_WRAP_GEMINI_EXIT") == "1" {
		wrapGemini([]string{"test"})
		return
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestWrapGemini_GeminiNotFound")
	cmd.Env = append(os.Environ(), "TEST_WRAP_GEMINI_EXIT=1", "PATH="+tmpDir)
	err := cmd.Run()

	exitErr := &exec.ExitError{}
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() != 1 {
			t.Errorf("Expected exit code 1, got %d", exitErr.ExitCode())
		}
	} else {
		t.Error("Expected command to exit with error")
	}
}

func TestRunGeminiDirectly_GeminiNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir)
	defer os.Setenv("PATH", oldPath)

	if os.Getenv("TEST_RUN_GEMINI_EXIT") == "1" {
		runGeminiDirectly([]string{"test"})
		return
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestRunGeminiDirectly_GeminiNotFound")
	cmd.Env = append(os.Environ(), "TEST_RUN_GEMINI_EXIT=1", "PATH="+tmpDir)
	err := cmd.Run()

	exitErr := &exec.ExitError{}
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() != 1 {
			t.Errorf("Expected exit code 1, got %d", exitErr.ExitCode())
		}
	} else {
		t.Error("Expected command to exit with error")
	}
}

func TestRunGeminiDirectly_Success(t *testing.T) {
	tmpDir := t.TempDir()

	mockGemini := filepath.Join(tmpDir, "gemini")
	mockScript := `#!/bin/bash
echo "Mock gemini output"
exit 0
`
	if err := os.WriteFile(mockGemini, []byte(mockScript), 0o755); err != nil {
		t.Fatal(err)
	}

	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir+":"+oldPath)
	defer os.Setenv("PATH", oldPath)

	if os.Getenv("TEST_RUN_GEMINI_SUCCESS") == "1" {
		runGeminiDirectly([]string{"test"})
		return
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestRunGeminiDirectly_Success")
	cmd.Env = append(os.Environ(), "TEST_RUN_GEMINI_SUCCESS=1", "PATH="+tmpDir+":"+oldPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected success, got error: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Mock gemini output") {
		t.Errorf("Expected mock output, got: %s", output)
	}
}

func TestRunGeminiDirectly_ExitCode(t *testing.T) {
	tmpDir := t.TempDir()

	mockGemini := filepath.Join(tmpDir, "gemini")
	mockScript := `#!/bin/bash
exit 42
`
	if err := os.WriteFile(mockGemini, []byte(mockScript), 0o755); err != nil {
		t.Fatal(err)
	}

	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir+":"+oldPath)
	defer os.Setenv("PATH", oldPath)

	if os.Getenv("TEST_RUN_GEMINI_EXITCODE") == "1" {
		runGeminiDirectly([]string{"test"})
		return
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestRunGeminiDirectly_ExitCode")
	cmd.Env = append(os.Environ(), "TEST_RUN_GEMINI_EXITCODE=1", "PATH="+tmpDir+":"+oldPath)
	err := cmd.Run()

	exitErr := &exec.ExitError{}
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() != 42 {
			t.Errorf("Expected exit code 42, got %d", exitErr.ExitCode())
		}
	} else {
		t.Error("Expected command to exit with error code 42")
	}
}
