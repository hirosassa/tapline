package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/hirosassa/tapline/pkg/logger"
	"github.com/hirosassa/tapline/pkg/session"
)

func wrapGemini(args []string) {
	sessionMgr, err := session.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: tapline session manager unavailable, logging disabled\n")
		runGeminiDirectly(args)
		return
	}

	log := logger.NewLogger("gemini-cli", sessionMgr)

	if !sessionMgr.HasActiveSession() {
		newSessionID := uuid.New().String()
		if err := sessionMgr.SetSessionID(newSessionID); err != nil {
			runGeminiDirectly(args)
			return
		}
		log.LogSessionStart(newSessionID, nil)
	}

	sessionID, err := sessionMgr.GetSessionID()
	if err != nil {
		runGeminiDirectly(args)
		return
	}

	if len(args) > 0 {
		prompt := strings.Join(args, " ")
		log.LogUserPrompt(sessionID, prompt)
	}

	response, exitCode := executeGemini(args)

	if response != "" {
		log.LogAssistantResponse(sessionID, response)
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func executeGemini(args []string) (response string, exitCode int) {
	geminiPath, err := exec.LookPath("gemini")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: 'gemini' command not found in PATH\n")
		return "", 1
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, geminiPath, args...)
	cmd.Stdin = os.Stdin

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdout pipe: %v\n", err)
		return "", 1
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stderr pipe: %v\n", err)
		return "", 1
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting gemini: %v\n", err)
		return "", 1
	}

	response = captureOutput(stdout, stderr)

	err = cmd.Wait()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
			return
		}
		exitCode = 1
		return
	}

	exitCode = 0
	return
}

func captureOutput(stdout, stderr io.Reader) string {
	lines := make(chan string, 100)
	stdoutDone := make(chan bool)
	stderrDone := make(chan bool)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			lines <- line
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdout: %v\n", err)
		}
		stdoutDone <- true
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintln(os.Stderr, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stderr: %v\n", err)
		}
		stderrDone <- true
	}()

	go func() {
		<-stdoutDone
		<-stderrDone
		close(lines)
	}()

	var response strings.Builder
	for line := range lines {
		response.WriteString(line)
		response.WriteString("\n")
	}

	return strings.TrimSpace(response.String())
}

func runGeminiDirectly(args []string) {
	geminiPath, err := exec.LookPath("gemini")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: 'gemini' command not found in PATH\n")
		os.Exit(1)
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, geminiPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}
