package secret

// Store fetches secrets needed for Bot integrations
type Store interface {
	TwitchToken() (string, error)
	ClientID() (string, error)
}
