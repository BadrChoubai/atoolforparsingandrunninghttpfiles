package cmd

import (
	"context"
	"github.com/spf13/cobra"

	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/vcs"
)

func NewRootCmd(ctx context.Context) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "atfparhf",
		Version: vcs.Version,
		Short:   "A tool for parsing and running http files",
	}

	rootCmd.AddCommand(newRunCmd(ctx))
	rootCmd.AddCommand(newVersionCmd())

	return rootCmd
}
