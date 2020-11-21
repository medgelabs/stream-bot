package bot

type (
	// Plugin that will send and receive messages like IRC
	Pluggable interface {
		Outbound
		Inbound
	}

	// Plugin will send messages to the bot e.g. PubSub
	Inbound interface {
		SetChannel(outbound chan<- Event)
	}

	// Plugin that will accept messages from the bot like Logger
	Outbound interface {
		GetChannel() chan<- Event
	}

	InboundPluginCollection  []Inbound
	OutboundPluginCollection []Outbound
)

// RegisterPlugin which both sends/receives messages to/from the Bot, respectively
func (bot *Bot) RegisterPlugin(plugin Pluggable) (err error) {
	err = bot.RegisterInboundPlugin(plugin)
	if err != nil {
		return err
	}

	err = bot.RegisterOutboundPlugin(plugin)
	if err != nil {
		return err
	}

	return
}

// RegisterInboundPlugin which sends messages TO the bot
func (bot *Bot) RegisterInboundPlugin(plugin Inbound) error {
	plugin.SetChannel(bot.events)
	bot.inboundPlugins = append(bot.inboundPlugins, plugin)
	return nil
}

// RegisterOutboundPlugin which receives messages FROM the bot
func (bot *Bot) RegisterOutboundPlugin(plugin Outbound) error {
	bot.outboundPlugins = append(bot.outboundPlugins, plugin)
	return nil
}
