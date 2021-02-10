package pubsub

import (
	"fmt"
	"medgebot/bot"
	"medgebot/ws/wstest"
	"testing"
	"time"
)

func TestMessageReceivedFromServer(t *testing.T) {
	conn := wstest.NewWebsocket()
	channelID := "testChannelID"
	authToken := "testAuthToken"

	client := NewClient(conn, channelID, authToken)
	testBot := make(chan bot.Event)
	client.SetDestination(testBot)
	client.Start()

	// Cursed be json-stringified `message` >.>
	conn.Send(fmt.Sprintf(`{
	  "type": "MESSAGE",
	  "data": {
		  "topic": "%s.%s",
		  "message": "{
			\"type\": \"reward-redeemed\",
			\"data\": {
			  \"timestamp\": \"2019-11-12T01:29:34.98329743Z\",
			  \"redemption\": {
				\"id\": \"9203c6f0-51b6-4d1d-a9ae-8eafdb0d6d47\",
				\"user\": {
				  \"id\": \"30515034\",
				  \"login\": \"testUser\",
				  \"display_name\": \"testUser\"
				},
				\"channel_id\": \"testChannelID\",
				\"redeemed_at\": \"2021-02-10T18:52:53.128421623Z\",
				\"reward\": {
				  \"id\": \"6ef17bb2-e5ae-432e-8b3f-5ac4dd774668\",
				  \"channel_id\": \"testChannelID\",
				  \"title\": \"Hydrate!\",
				  \"prompt\": \"\",
				  \"cost\": 10,
				  \"is_user_input_required\": false,
				  \"is_sub_only\": false,
				  \"image\": {
					\"url_1x\": \"https://static-cdn.jtvnw.net/custom-reward-images/30515034/6ef17bb2-e5ae-432e-8b3f-5ac4dd774668/7bcd9ca8-da17-42c9-800a-2f08832e5d4b/custom-1.png\",
					\"url_2x\": \"https://static-cdn.jtvnw.net/custom-reward-images/30515034/6ef17bb2-e5ae-432e-8b3f-5ac4dd774668/7bcd9ca8-da17-42c9-800a-2f08832e5d4b/custom-2.png\",
					\"url_4x\": \"https://static-cdn.jtvnw.net/custom-reward-images/30515034/6ef17bb2-e5ae-432e-8b3f-5ac4dd774668/7bcd9ca8-da17-42c9-800a-2f08832e5d4b/custom-4.png\"
				  },
				  \"default_image\": {
					\"url_1x\": \"https://static-cdn.jtvnw.net/custom-reward-images/default-1.png\",
					\"url_2x\": \"https://static-cdn.jtvnw.net/custom-reward-images/default-2.png\",
					\"url_4x\": \"https://static-cdn.jtvnw.net/custom-reward-images/default-4.png\"
				  },
				  \"background_color\": \"#00C7AC\",
				  \"is_enabled\": true,
				  \"is_paused\": false,
				  \"is_in_stock\": true,
				  \"max_per_stream\": { \"is_enabled\": false, \"max_per_stream\": 0 },
				  \"should_redemptions_skip_request_queue\": true
				},
				\"user_input\": \"\",
				\"status\": \"FULFILLED\"
			  }
			}
		 }"
	  }`, ChannelPointTopic, channelID))

	// Wait for message on bot Event channel
	select {
	case evt := <-testBot:
		if evt.Type != bot.POINT_REDEMPTION {
			t.Fatalf("Did not receive a ChannelPoint message. Got %+v", evt)
		}

		if evt.Amount != 10 {
			t.Fatalf("Did not receive the correct event Amount. Got %d", evt.Amount)
		}

		if evt.Title != "Hydrate!" {
			t.Fatalf("Did not receive the correct event Title. Got %s", evt.Title)
		}
	case <-time.After(3 * time.Second):
		fmt.Println(conn.String())
		t.Fatalf("Timeout while waiting to receive expected message")
	}
}
