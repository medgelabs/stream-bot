package bot

import (
	"text/template"
)

func makeTemplate(name, tmpl string) *template.Template {
	return template.Must(template.New(name).Parse(tmpl))
}

// Register assertions as inboundPlugin
type Checker struct {
	events chan Event
}

func NewChecker() *Checker {
	return &Checker{
		events: make(chan Event),
	}
}

func (c *Checker) GetChannel() chan<- Event {
	return c.events
}
