package greeter

// Config for a greeter in a channel
type Config struct {
	// MessageFormat is a go fmt-style formatted string
	MessageFormat string `mapstructure:"messageFormat"`
}
