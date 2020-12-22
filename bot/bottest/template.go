package bottest

import "text/template"

// Helper for creating templates for use in Handlers
func MakeTemplate(name, tmpl string) *template.Template {
	return template.Must(template.New(name).Parse(tmpl))
}
