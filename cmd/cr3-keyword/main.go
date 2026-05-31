package main

import (
	"context"
	"fmt"
	"os"

	"cr3-keywords/internal/cli"
	"cr3-keywords/internal/cli/shell"
	"cr3-keywords/internal/logging"
	"cr3-keywords/internal/pipeline"
)

var version = "dev"
var buildTime = "unknown"

func main() {
	opts, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if opts.Version {
		fmt.Fprintf(os.Stdout, "%s (%s)\n", version, buildTime)
		return
	}

	if opts.Help {
		fmt.Print(cli.Usage())
		return
	}

	if opts.Clear {
		deleted, err := cli.ClearMatchedFiles(opts.CR3Path, opts.Files)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Deleted %d file(s).\n", deleted)
		return
	}

	if opts.Install {
		installedPath, err := shell.Install(opts.InstallTarget, "cr3")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Installed %s completion/module: %s\n", opts.InstallTarget, installedPath)
		fmt.Fprintf(os.Stdout, "%s\n", installHint(opts.InstallTarget, installedPath))
		return
	}

	if opts.Completion {
		script, err := shell.Script(opts.CompletionTarget, "cr3")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, script)
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

func installHint(target, path string) string {
	switch target {
	case "nushell":
		return fmt.Sprintf("Restart Nushell or run: source %s", path)
	case "zsh":
		return fmt.Sprintf("Run: autoload -Uz compinit && compinit (ensure completion dir is in fpath, then restart zsh). Installed file: %s", path)
	case "bash":
		return fmt.Sprintf("Load in current shell: source %s (or restart bash).", path)
	default:
		return ""
	}
}
