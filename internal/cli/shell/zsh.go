package shell

func renderZshCompletion(commandName string) (string, error) {
	return renderTemplate("zsh.tmpl.zsh", commandName)
}
