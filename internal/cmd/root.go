package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/lucasepe/checkit/internal/cmd/render"
	getoptutil "github.com/lucasepe/checkit/internal/util/getopt"
	ioutil "github.com/lucasepe/checkit/internal/util/io"

	"github.com/lucasepe/x/getopt"
)

var (
	BuildKey = buildKey{}
)

type Action int

const (
	NoAction Action = iota
	Render
	ShowHelp
	ShowVersion
)

func Run(ctx context.Context) (err error) {
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

	err = render.Do(os.Args[1:])
	if errors.Is(err, ioutil.ErrNoInputDetected) {
		usage(os.Stderr)
		return nil
	}

	return err
}

func chosenAction(args []string) (Action, error) {
	_, opts, err := getopt.GetOpt(args,
		"hv",
		[]string{"help", "version"},
	)
	if err != nil {
		return NoAction, err
	}

	showVersion := getoptutil.HasOpt(opts, []string{"-v", "--version"})
	if showVersion {
		return ShowVersion, nil
	}

	showHelp := getoptutil.HasOpt(opts, []string{"-h", "--help"})
	if showHelp {
		return ShowHelp, nil
	}

	return Render, nil
}

type (
	buildKey struct{}
)
