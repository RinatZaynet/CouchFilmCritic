package main

import "text/template"

func ParseTemplates(pattern string) *template.Template {
	// Сделать нормальную обработку ошибок
	return template.Must(template.ParseGlob(pattern))
}
