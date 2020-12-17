package bot

import (
	"text/template"
)

// Register assertions as inboundPlugin
type Checker struct {
	Events chan Event
}

func NewChecker() *Checker {
	return &Checker{
		Events: make(chan Event),
	}
}

func (c *Checker) GetChannel() chan<- Event {
	return c.Events
}

// Helper for creating templates for use in Handlers
func makeTemplate(name, tmpl string) *template.Template {
	return template.Must(template.New(name).Parse(tmpl))
}
