package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/parser"
)

func newRunCmd(ctx context.Context) *cobra.Command {
	var (
		filename string
		curl     bool
	)

	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse a .http_file file",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
			defer cancel()
			return run(ctx, os.Stdin, os.Stdout, os.Stderr, filename, curl)
		},
	}

	cmd.Flags().StringVarP(&filename, "file", "f", "", "Path to the .http_file")
	cmd.Flags().BoolVar(&curl, "curl", false, "Generate curl output")
	cmd.MarkFlagRequired("file")

	return cmd
}

func run(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, filename string, curl bool) error {
	httpFileParser := parser.NewHttpFileParser()

	parsed, err := httpFileParser.Parse(filepath.Clean(filename))
	if err != nil {
		return fmt.Errorf("failed to parse request file: %w", err)
	}
	if !parsed || len(httpFileParser.ScannedLines) == 0 {
		return nil
	}

	requests, err := httpFileParser.BuildRequests()
	if err != nil {
		return err
	}

	if curl {
		file, errCreateFile := os.Create("generated_curl.sh")
		if errCreateFile != nil {
			return fmt.Errorf("error creating file: %w", errCreateFile)
		}
		defer file.Close()

		file.WriteString("#!/bin/bash\n")

		for i, request := range requests {
			c, err := request.ToCurl()
			if err != nil {
				return err
			}

			if i != 0 {
				file.WriteString("\n")
			}
			fmt.Fprintf(file, "( # %s\n\t%s\n)\n", request.Description, c)
		}
	}

	return nil
}
