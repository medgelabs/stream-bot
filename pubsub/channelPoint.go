package pubsub

import "time"

// ChannelPointRedemption represents a ChannelPoint message from PubSub,
// cut down to the data we need
type ChannelPointRedemption struct {
	Type string `json:"type"`
	Data struct {
		Redemption struct {
			ID   string `json:"id"`
			User struct {
				ID          string `json:"id"`
				Login       string `json:"login"`
				DisplayName string `json:"display_name"`
			} `json:"user"`
			ChannelID  string    `json:"channel_id"`
			RedeemedAt time.Time `json:"redeemed_at"`
			Reward     struct {
				ID        string `json:"id"`
				Title     string `json:"title"`
				Prompt    string `json:"prompt"`
				Cost      int    `json:"cost"`
				IsEnabled bool   `json:"is_enabled"`
				IsPaused  bool   `json:"is_paused"`
				IsInStock bool   `json:"is_in_stock"`
			} `json:"reward"`
			UserInput string `json:"user_input"`
		} `json:"redemption"`
	} `json:"data"`
}
