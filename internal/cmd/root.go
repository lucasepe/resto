package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/lucasepe/resto/internal/cmd/call"
	"github.com/lucasepe/resto/internal/env"
	getoptutil "github.com/lucasepe/resto/internal/util/getopt"
	ioutil "github.com/lucasepe/resto/internal/util/io"

	"github.com/lucasepe/x/getopt"
)

var (
	BuildKey = buildKey{}
)

type Action int

const (
	NoAction Action = iota
	Call
	ShowHelp
	ShowVersion
)

func Run(ctx context.Context) (err error) {
	file, err := env.DefaultEnvFile()
	if err != nil {
		return err
	}

	err = env.Load(file, true)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	act, err := chosenAction(os.Args[1:])
	switch act {
	case ShowHelp:
		usage(os.Stderr)
		return nil
	case ShowVersion:
		bld := ctx.Value(BuildKey).(string)
		fmt.Fprintf(os.Stderr, "%s - build: %s\n", appName, bld)
		return nil
	}

	err = call.Do(os.Args[1:])
	if errors.Is(err, ioutil.ErrNoInputDetected) {
		usage(os.Stderr)
		return nil
	}

	return err
}

func chosenAction(args []string) (Action, error) {
	_, opts, err := getopt.GetOpt(args,
		"",
		[]string{"help", "version"},
	)
	if err != nil {
		return NoAction, err
	}

	showVersion := getoptutil.HasOpt(opts, []string{"--version"})
	if showVersion {
		return ShowVersion, nil
	}

	showHelp := getoptutil.HasOpt(opts, []string{"--help"})
	if showHelp {
		return ShowHelp, nil
	}

	return Call, nil
}

type (
	buildKey struct{}
)
