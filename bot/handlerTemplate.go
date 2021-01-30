package bot

import (
	"log"
	"strings"
	"text/template"
)

type HandlerTemplate struct {
	template *template.Template
}

func NewHandlerTemplate(t *template.Template) HandlerTemplate {
	return HandlerTemplate{
		template: t,
	}
}

func (h HandlerTemplate) Parse(evt Event) string {
	var msg strings.Builder
	err := h.template.Execute(&msg, evt)
	if err != nil {
		log.Printf("ERROR: bits template execute - %v", err)
		return "" // We assume bot will not send empty messages
	}

	return msg.String()
}
