package secret

// Store fetches secrets needed for Bot integrations
type Store interface {
	GetTwitchToken() (string, error)
}
