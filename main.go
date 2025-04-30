package main

import (
	"context"
	"os"

	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/cmd"
)

func main() {
	ctx := context.Background()
	if err := cmd.NewRootCmd(ctx).Execute(); err != nil {
		os.Exit(1)
	}
}
