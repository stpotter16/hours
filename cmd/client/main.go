package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"text/tabwriter"

	"github.com/stpotter16/hours/internal/client"
	"github.com/stpotter16/hours/internal/config"
	"github.com/stpotter16/hours/internal/types"
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

	if len(args) < 2 {
		return errors.New("Usage: hours <command>")
	}

	cmd := args[1]

	// configure doesn't need a client
	if cmd == "configure" {
		return runConfigure(stdout)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if host := getenv("HOURS_HOST"); host != "" {
		cfg.Host = host
	}

	c, err := client.New(cfg.Host)
	if err != nil {
		return err
	}

	switch cmd {
	case "login":
		passphrase := cfg.Passphrase
		if passphrase == "" {
			fmt.Fprint(stdout, "Passphrase: ")
			pb, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			fmt.Fprintln(stdout)
			passphrase = string(pb)
		}
		if err := c.Login(ctx, passphrase); err != nil {
			return err
		}
	case "projects":
		if len(args) < 3 {
			return errors.New("Usage: hours projects <list|create|delete>")
		}
		switch args[2] {
		case "list":
			projects, err := c.ListProjects(ctx)
			if err != nil {
				return err
			}
			printProjects(stdout, projects)
		case "create":
			if len(args) < 4 {
				return errors.New("Usage: hours projects create <name>")
			}
			if err := c.CreateProject(ctx, args[3]); err != nil {
				return err
			}
		case "delete":
			if len(args) < 4 {
				return errors.New("Usage: hours projects delete <name>")
			}
			if err := c.DeleteProject(ctx, args[3]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Invalid subcommand: %s", args[2])
		}
	case "timers":
		if len(args) < 3 {
			return errors.New("Usage: hours timers <list|start|stop>")
		}
		switch args[2] {
		case "list":
			timers, err := c.ListTimers(ctx)
			if err != nil {
				return err
			}
			printTimers(stdout, timers)
		case "start":
			if len(args) < 4 {
				return errors.New("Usage: hours timers start <project name>")
			}
			if err := c.StartTimer(ctx, args[3]); err != nil {
				return err
			}
		case "stop":
			if len(args) < 4 {
				return errors.New("Usage: hours timers stop <project name>")
			}
			if err := c.StopTimer(ctx, args[3]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Invalid subcommand: %s", args[2])
		}
	default:
		return fmt.Errorf("Invalid command: %s", cmd)
	}

	return nil
}

func runConfigure(stdout io.Writer) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Fprint(stdout, "Host URL: ")
	var host string
	fmt.Scanln(&host)
	if host != "" {
		cfg.Host = host
	}

	fmt.Fprint(stdout, "Passphrase (leave blank to skip): ")
	passphrase, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout)
	if len(passphrase) > 0 {
		cfg.Passphrase = string(passphrase)
	}

	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Fprintln(stdout, "Configuration saved.")
	return nil
}

func printProjects(w io.Writer, resp types.ProjectListResponse) {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "ID\tNAME\tTIME WORKED\tCREATED\tLAST MODIFIED")
	for _, p := range resp.Projects {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\n",
			p.ID,
			p.Name,
			formatDuration(p.TotalSeconds),
			p.CreatedTime.Format("2006-01-02"),
			p.LastModifiedTime.Format("2006-01-02"),
		)
	}
	tw.Flush()
}

func formatDuration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	if h > 0 {
		return fmt.Sprintf("%dh %02dm %02ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func printTimers(w io.Writer, resp types.TimerListResponse) {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "ID\tPROJECT NAME\tSTARTED")
	for _, t := range resp.Timers {
		fmt.Fprintf(tw, "%d\t%s\t%s\n",
			t.ID,
			t.ProjectName,
			t.StartedTime.Format("2006-01-02"),
		)
	}
	tw.Flush()
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
