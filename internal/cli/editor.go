package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EditPromptInTerminal ensures the prompt file exists and opens it in the
// user's terminal editor.
//
// The editor is resolved from the EDITOR environment variable, and defaults
// to "nano" when EDITOR is not set.
func EditPromptInTerminal(promptPath string) error {
	if err := ensurePromptFile(promptPath); err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	cmd, err := editorCommand(editor, promptPath)
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("open editor %q: %w", editor, err)
	}

	return nil
}

func ensurePromptFile(promptPath string) error {
	if _, err := os.Stat(promptPath); err == nil {
		// does not check if the file is empty
		return nil
	} else if !os.IsNotExist(err) {
		// if the real error is not "file does not exist", return it
		return err
	}

	if err := os.MkdirAll(filepath.Dir(promptPath), 0o755); err != nil {
		return err
	}

	defaultPrompt := `Context:
This photo is part of a series of photos.

Task:
1. Describe the image
2. Generate 10–20 simple keywords (comma-separated)
3. Write a short caption (1 sentence), using the image description from (task 1).

Avoid generic terms like "image" or "photo".

Output should be result of task 2, empty line, result of task 3.
Make sure output does not contain the full description,
only comma separated keywords in the first line, an empty line and the caption.
`

	return os.WriteFile(promptPath, []byte(defaultPrompt), 0o644)
}

func editorCommand(editor, promptPath string) (*exec.Cmd, error) {
	parts := strings.Fields(editor)
	if len(parts) == 0 {
		return nil, fmt.Errorf("EDITOR is empty")
	}

	args := append(parts[1:], promptPath)
	return exec.Command(parts[0], args...), nil
}
