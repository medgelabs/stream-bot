package bot

import "errors"

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

// This will register a plugin which is an inbound and outbound plugin
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

// This will register an inbound plugin
func (bot *Bot) RegisterInboundPlugin(plugin Inbound) error {
	_, ok := bot.inboundPlugins[plugin.GetId()]
	if ok {
		return errors.New("inbound plugin already registered")
	}
	plugin.SetOutboundChannel(bot.events)
	bot.inboundPlugins[plugin.GetId()] = plugin

	return nil
}

// This will register an outbound plugin
func (bot *Bot) RegisterOutboundPlugin(plugin Outbound) error {
	_, ok := bot.outboundPlugins[plugin.GetId()]
	if ok {
		return errors.New("outbound plugin already registered")
	}

	bot.outboundPlugins[plugin.GetId()] = plugin

	return nil
}
