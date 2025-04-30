package main

import (
	"context"
	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/cmd"
	"os"
)

func main() {
	ctx := context.Background()
	if err := cmd.NewRootCmd(ctx).Execute(); err != nil {
		os.Exit(1)
	}
}
