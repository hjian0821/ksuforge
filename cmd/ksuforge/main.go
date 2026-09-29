package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/hjian0821/ksuforge/internal/cli"
	"github.com/hjian0821/ksuforge/internal/logging"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, closeLog := logging.Setup(logging.Options{Level: slog.LevelInfo, Prefix: "ksuforge"})
	defer closeLog()
	logger.Info("app.started", "args", os.Args[1:])

	code := cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	if code != 0 {
		logger.Error("app.command_failed", "exit_code", code, "args", os.Args[1:])
	}
	os.Exit(code)
}
