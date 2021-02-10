package bot

// Client represents data being sent TO the Bot
type Client interface {
	SetDestination(chan<- Event)
}
