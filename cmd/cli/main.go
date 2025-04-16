package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/logging"
	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/parser"
)

type config struct {
	filename string
}

func run(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	logger := logging.NewLogger(stdout, stderr)

	var cfg config

	flag.StringVar(&cfg.filename, "file", "", "http file")

	flag.Parse()

	httpFileParser := &parser.HttpFileParser{}

	logger.Info("parsing .http file")
	parsed, err := httpFileParser.Parse(filepath.Clean(cfg.filename))
	if err != nil {
		logger.Error("parsing http file", "error", err)
		return err
	}

	if parsed {
		logger.Info(
			"parsed .http file",
			"requests parsed from file",
			len(httpFileParser.Requests))
	}

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
