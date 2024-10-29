package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/RIDOS/sub-cmd/cmd"
)

var totalDuration time.Duration = 10
var errInvalidSubCommand = errors.New("Invalid sub-command specified")

func printUsage(ctx context.Context,w io.Writer) {
	fmt.Fprintf(w, "Usage: mync [http|grpc] -h\n")
	cmd.HandleHttp(ctx, w, []string{"-h"})
	cmd.HandleGrpc(w, []string{"-h"})
}

func handleCommand(ctx context.Context, w io.Writer, args []string) error {
	var err error

	if len(args) < 1 {
		err = errInvalidSubCommand
	} else {
		switch args[0] {
		case "http":
			err = cmd.HandleHttp(ctx, w, args[1:])
		case "grpc":
			err = cmd.HandleGrpc(w, args[1:])
		case "-h":
			printUsage(ctx, w)
		case "--help":
			printUsage(ctx, w)
		default:
			err = errInvalidSubCommand
		}
	}

	if errors.Is(err, cmd.ErrNoServerSpecified) || errors.Is(err, errInvalidSubCommand) || errors.Is(err, cmd.ErrInvalidMethod) {
		fmt.Fprintln(w, err)
		printUsage(ctx, w)
	}

	return err
}

func main() {
	allowedDuration := totalDuration * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), allowedDuration)
	defer cancel()

	chanel := make(chan error, 1)

	go func() {
		err := handleCommand(ctx, os.Stdout, os.Args[1:])
		chanel <- err
	}()

	select {
	case <-ctx.Done():
		fmt.Fprintln(os.Stdout, "Time out close...")
		os.Exit(0)
	case err := <-chanel:
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
