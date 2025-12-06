package main

import (
	"context"
	"errors"
	"io"
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

func TestCaptureOutput_BasicOutput(t *testing.T) {
	stdout := strings.NewReader("line 1\nline 2\nline 3")
	stderr := strings.NewReader("")

	result := captureOutput(stdout, stderr)

	expected := "line 1\nline 2\nline 3"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCaptureOutput_EmptyOutput(t *testing.T) {
	stdout := strings.NewReader("")
	stderr := strings.NewReader("")

	result := captureOutput(stdout, stderr)

	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestCaptureOutput_WithTrailingNewline(t *testing.T) {
	stdout := strings.NewReader("line 1\nline 2\n")
	stderr := strings.NewReader("")

	result := captureOutput(stdout, stderr)

	expected := "line 1\nline 2"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCaptureOutput_StderrIgnored(t *testing.T) {
	stdout := strings.NewReader("stdout line")
	stderr := strings.NewReader("stderr line")

	result := captureOutput(stdout, stderr)

	expected := "stdout line"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCaptureOutput_MultipleLines(t *testing.T) {
	lines := strings.Repeat("line\n", 100)
	stdout := strings.NewReader(lines)
	stderr := strings.NewReader("")

	result := captureOutput(stdout, stderr)

	if !strings.Contains(result, "line") {
		t.Error("Expected output to contain 'line'")
	}

	lineCount := strings.Count(result, "\n") + 1
	if lineCount != 100 {
		t.Errorf("Expected 100 lines, got %d", lineCount)
	}
}

func TestCaptureOutput_ConcurrentReads(t *testing.T) {
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()

	go func() {
		w1.Write([]byte("line 1\n"))
		w1.Write([]byte("line 2\n"))
		w1.Close()
	}()

	go func() {
		w2.Write([]byte("error 1\n"))
		w2.Close()
	}()

	result := captureOutput(r1, r2)

	if !strings.Contains(result, "line 1") {
		t.Error("Expected output to contain 'line 1'")
	}
	if !strings.Contains(result, "line 2") {
		t.Error("Expected output to contain 'line 2'")
	}
}
