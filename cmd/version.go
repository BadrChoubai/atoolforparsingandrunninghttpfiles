package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/vcs"
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
