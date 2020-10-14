package secret

// StaticStore is an in-memory secrets store. Don't use this for Prod!
type StaticStore struct {
	twitchToken string
}

func NewStaticStore(twitchToken string) StaticStore {
	return StaticStore{
		twitchToken: twitchToken,
	}
}

func (s StaticStore) GetTwitchToken() (string, error) {
	return s.twitchToken, nil
}
