package cmd

import (
	"fmt"
	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/vcs"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version:", vcs.Version)
		},
	}
}
