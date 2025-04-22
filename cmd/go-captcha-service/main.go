/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

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
	fmt.Fprintf(os.Stdout, "[Main] Starting the application ...")
	a, err := app.NewApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Main] Failed to initialize app: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start app
	if err = a.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "[Main] Failed to start app: %v\n", err)
	}

	// Handle termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	fmt.Fprintf(os.Stdout, "[Main] The application start successfully")
	<-sigCh

	a.Shutdown()
	fmt.Fprintf(os.Stderr, "[Main] App service exited")
}
