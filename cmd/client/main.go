package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/stpotter16/hours/internal/client"
	"golang.org/x/term"
)

func run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := client.New(getenv)
	if err != nil {
		return err
	}

	if len(args) < 2 {
		return errors.New("Invalid number of args")
	}

	cmd := args[1]

	switch cmd {
	case "login":
		fmt.Print("Passphrase:")
		passphrase, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Println()
		err = client.Login(string(passphrase))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid command: %s", cmd)
	}

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(
		ctx,
		os.Args,
		os.Getenv,
		nil,
		os.Stdout,
		os.Stderr,
	); err != nil {
		log.Fatalf("%s", err)
	}
}
