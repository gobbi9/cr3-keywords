package cli

import (
	"fmt"
	"os"
	"os/exec"
)

func EditPromptInTerminal(promptPath string) error {
	if err := ensurePromptFile(promptPath); err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, promptPath)
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
		return nil
	}

	initial := `Context:
This photo is part of a series of photos taken in Braunschweig, Germany.

Task:
1. Describe the image
2. Generate 10–20 simple keywords (comma-separated)
3. Write a short caption (1 sentence)

Output should be result of task 2, empty line, result of task 3.
`

	return os.WriteFile(promptPath, []byte(initial), 0o644)
}
