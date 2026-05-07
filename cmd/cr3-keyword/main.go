package main

import (
	"context"
	"fmt"
	"os"

	"cr3-keywords/internal/cli"
	"cr3-keywords/internal/logging"
	"cr3-keywords/internal/pipeline"
)

func main() {
	opts, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if opts.Help {
		fmt.Print(cli.Usage())
		return
	}

	if opts.Clear {
		deleted, err := cli.ClearLastRunXMP(os.Stdin, os.Stdout)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Deleted %d .xmp file(s).\n", deleted)
		return
	}

	logger := logging.New(opts.Verbose)

	if opts.EditPrompt {
		if err := cli.EditPromptInTerminal(opts.PromptPath); err != nil {
			logger.Error("failed to open prompt editor", "error", err)
			os.Exit(1)
		}
	}

	runner := pipeline.NewRunner(opts, logger)
	if err := runner.Run(context.Background()); err != nil {
		logger.Error("pipeline failed", "error", err)
		os.Exit(1)
	}

	if !opts.DryRun {
		if err := cli.SaveLastRunCR3Path(opts.CR3Path); err != nil {
			logger.Warn("failed to save last run CR3 path", "error", err)
		}
	}
}
