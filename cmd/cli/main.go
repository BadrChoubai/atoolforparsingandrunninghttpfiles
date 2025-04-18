package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/httpclient"
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

	if cfg.filename == "" {
		return fmt.Errorf("no input file specified. Use --file to provide a .http file")
	}

	httpFileParser := &parser.HttpFileParser{}

	logger.Info("parsing requests in file", "file", cfg.filename)

	parsed, err := httpFileParser.Parse(filepath.Clean(cfg.filename))
	if err != nil {
		return fmt.Errorf("failed to parse request file: %w", err)
	}
	if !parsed || len(httpFileParser.Requests) == 0 {
		logger.Info("no valid HTTP requests found in file")
		return nil
	}

	logger.Info(
		"parsed .http file",
		"requests parsed from file",
		len(httpFileParser.Requests))

	httpClient := httpclient.NewHTTPClient(logger)
	for _, r := range httpFileParser.Requests {
		logger.Info(
			"sending request",
			"method",
			r.Method,
			"url",
			r.URL,
		)
		response, httpClientError := httpClient.DoRequest(r)
		if httpClientError != nil {
			logger.Error("trying request", "error", httpClientError.Error(), "method", r.Method, "url", r.URL)
			return err
		}

		logger.Info(
			"response",
			"status", response.Status)
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
