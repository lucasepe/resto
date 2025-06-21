package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/lucasepe/resto/internal/cmd"
)

var (
	Version = "v0.0.0"
)

func main() {
	log.SetOutput(os.Stderr)

	ctx := context.Background()
	ctx = context.WithValue(ctx, cmd.BuildKey, Version)

	if err := cmd.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
