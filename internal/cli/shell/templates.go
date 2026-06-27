package shell

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed *.tmpl.sh *.tmpl.zsh *.tmpl.nu
var shellTemplates embed.FS

type templateData struct {
	CommandName string
}

func renderTemplate(templateName, commandName string) (string, error) {
	tpl, err := template.ParseFS(shellTemplates, templateName)
	if err != nil {
		return "", fmt.Errorf("parse template %s: %w", templateName, err)
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, templateData{CommandName: commandName}); err != nil {
		return "", fmt.Errorf("render template %s: %w", templateName, err)
	}
	return out.String(), nil
}
