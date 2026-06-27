package shell

func renderBashCompletion(commandName string) (string, error) {
	return renderTemplate("bash.tmpl.sh", commandName)
}
