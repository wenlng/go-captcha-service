package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wenlng/go-captcha-service/internal/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize app: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = a.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start app: %v\n", err)
	}

	// Handle termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	a.Shutdown()
	fmt.Fprintf(os.Stderr, "App service exited")
}
