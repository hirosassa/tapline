package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
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

	geminiPath, err := exec.LookPath("gemini")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: 'gemini' command not found in PATH\n")
		os.Exit(1)
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, geminiPath, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdout pipe: %v\n", err)
		os.Exit(1)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stderr pipe: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting gemini: %v\n", err)
		os.Exit(1)
	}

	var response strings.Builder

	stdoutDone := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			response.WriteString(line)
			response.WriteString("\n")
		}
		stdoutDone <- true
	}()

	stderrDone := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintln(os.Stderr, scanner.Text())
		}
		stderrDone <- true
	}()

	<-stdoutDone
	<-stderrDone

	err = cmd.Wait()

	if response.Len() > 0 {
		log.LogAssistantResponse(sessionID, strings.TrimSpace(response.String()))
	}

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
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
