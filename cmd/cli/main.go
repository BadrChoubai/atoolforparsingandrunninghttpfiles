package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/logging"
	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/parser"
	"io"
	"os"
	"os/signal"
	"path/filepath"
)

var (
	filename string
)

func run(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	logger := logging.NewLogger(stdout, stderr)

	flag.StringVar(&filename, "file", "", "http_file file")
	flag.Parse()

	if filename == "" {
		return fmt.Errorf("no input file specified. Use --file to provide a .http_file file")
	}

	httpFileParser := parser.NewHttpFileParser()

	logger.Info("parsing requests in file", "file", filename)

	parsed, err := httpFileParser.Parse(filepath.Clean(filename))
	if err != nil {
		return fmt.Errorf("failed to parse request file: %w", err)
	}
	if !parsed || len(httpFileParser.ScannedLines) == 0 {
		return nil
	}

	logger.Info(
		"parsed .http_file file",
		"requests parsed from file",
		len(httpFileParser.ScannedLines))

	requests, err := httpFileParser.BuildRequests()
	fmt.Println(requests[1].Request.URL.String())

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
