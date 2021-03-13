package bot

import (
	log "medgebot/logger"
	"strings"
	"text/template"
)

// HandlerTemplate wraps text/template for convenience
type HandlerTemplate struct {
	template *template.Template
}

// NewHandlerTemplate creates a new instance. Intent is to call Parse() soon after
func NewHandlerTemplate(t *template.Template) HandlerTemplate {
	return HandlerTemplate{
		template: t,
	}
}

// Parse interpolates the given Event onto the stored template
func (h HandlerTemplate) Parse(evt Event) string {
	var msg strings.Builder
	err := h.template.Execute(&msg, evt)
	if err != nil {
		log.Error(err, "template execute")
		return "" // We assume bot will not send empty messages
	}

	return msg.String()
}
