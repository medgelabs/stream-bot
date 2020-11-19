package bot

type (
	// Plugin that will send and receive messages like IRC
	Pluggable interface {
		Outbound
		Inbound
	}

	// Plugin will send messages to the bot e.g. PubSub
	Inbound interface {
		GetId() string                            // for uniqueness check
		SetOutboundChannel(outbound chan<- Event) // channel that is used by the plugin
	}

	// Plugin that will accept messages from the bot like Logger
	Outbound interface {
		GetId() string                   // for uniqueness check
		GetInboundChannel() chan<- Event // channel that is used by the bot
	}

	InboundPluginCollection  map[string]Inbound
	OutboundPluginCollection map[string]Outbound
)
