package app

import (
	"errors"
	"html/template"
	"os"
)

var (
	errTemplatePathNotExist = errors.New("template path does not exist")
	errInvalidTemplatePath  = errors.New("invalid template path")
)

func parseTemplates(templatePath string) (*template.Template, error) {
	if templatePath == "" {
		return nil, errTemplatePathNotExist
	}
	if _, err := os.Stat(templatePath); err != nil {
		return nil, errInvalidTemplatePath
	}
	tmpl, err := template.ParseGlob(templatePath + "*.html")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
