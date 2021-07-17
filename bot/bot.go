package bot

import (
	"errors"
	"fmt"
	"medgebot/cache"
	"medgebot/logger"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Bot handles various feature processing based on Stream events
type Bot struct {
	sync.Mutex

	// Handlers of bot.Event messages
	consumers []Handler

	// Producers of data from external systems
	clients []Client

	// Where the Bot sends messages to get to Chat
	chatClient ChatClient

	// events
	events    chan Event
	listening bool

	// Cache for various handler metrics
	dataStore cache.Cache

	// polls
	pollRunning  bool
	pollQuestion string
	pollAnswers  []string
}

// New produces a newly instantiated Bot
func New(metricsCache cache.Cache) Bot {
	return Bot{
		consumers:   make([]Handler, 0),
		clients:     make([]Client, 0),
		events:      make(chan Event, 0),
		listening:   false,
		dataStore:   metricsCache,
		pollRunning: false,
	}
}

// Start the bot and listen for incoming events
func (bot *Bot) Start() error {
	// Ensure single concurrent reader, per doc requirements
	go bot.listen()

	return nil
}

// RegisterClient links a Client that will send data TO the Bot.
// This method also set's the Client's Destination channel
func (bot *Bot) RegisterClient(client Client) {
	bot.Lock()
	defer bot.Unlock()

	client.SetDestination(bot.events)
	bot.clients = append(bot.clients, client)
}

// SetChatClient register's the client that will allow the Bot to send
// messages to Chat
func (bot *Bot) SetChatClient(client ChatClient) {
	bot.Lock()
	defer bot.Unlock()

	bot.chatClient = client
}

// RegisterHandler registers a function that will be called concurrently when a message is received
func (bot *Bot) RegisterHandler(consumer Handler) error {
	if bot.listening {
		return errors.New("RegisterHandler called after bot already listening")
	}

	bot.Lock()
	defer bot.Unlock()

	consumers := append(bot.consumers, consumer)
	bot.consumers = consumers
	return nil
}

// ReceiveEvent is a way for code to directly queue Events to be processed. Ex: alias commands
func (bot *Bot) ReceiveEvent(evt Event) {
	bot.events <- evt
}

// SendMessage sends a message to the given channel, without prefix
func (bot *Bot) SendMessage(message string, args ...interface{}) {
	if strings.TrimSpace(message) == "" {
		return
	}

	evt := NewChatEvent()
	evt.Message = fmt.Sprintf(message, args...)

	go bot.sendEvent(evt)
}

// IsPollRunning checks if a Poll is currently active
func (bot *Bot) IsPollRunning() bool {
	return bot.pollRunning
}

// StartPoll starts a poll within the Bot. Returns error if poll already running.
// Store reference to Poll contents for Handler use later
func (bot *Bot) StartPoll(duration time.Duration, question string, answers []string) error {
	if bot.IsPollRunning() {
		return errors.New("Poll already running")
	}

	logger.Info("Starting new Poll: %s", question)
	bot.pollRunning = true
	bot.pollQuestion = question
	bot.pollAnswers = answers

	bot.SendPollMessage()

	// Spawn off goroutine to close the poll after the given Duration
	go func(bot *Bot, dur time.Duration) {
		select {
		case <-time.After(dur):
			bot.closePoll()
		}
	}(bot, duration)

	return nil
}

// SendPollMessage sends the current Poll message, if a poll is running
func (bot *Bot) SendPollMessage() {
	if !bot.pollRunning {
		bot.SendMessage("No poll running")
		return
	}

	formattedAnswers := ""
	for idx, answer := range bot.pollAnswers {
		formattedAnswers += fmt.Sprintf("%d: %s | ", idx+1, answer)
	}

	bot.SendMessage("Poll started! Question: %s : %s", bot.pollQuestion, formattedAnswers)
}

// closePoll ends an active poll
func (bot *Bot) closePoll() {
	// Count each vote
	answersStr, err := bot.dataStore.Get("pollAnswers")
	if err != nil {
		logger.Error(err, "Failed to fetch pollAnswers to close poll")
		return
	}
	splitAnswers := strings.Split(answersStr, ",")

	answerCounts := make([]int, len(bot.pollAnswers))
	for _, answerStr := range splitAnswers {
		answer, err := strconv.Atoi(answersStr)
		if err != nil {
			logger.Warn("Answer %s is not a valid number. Skipping", answerStr)
			continue
		}

		answerCounts[answer-1]++
	}

	highestIdx := 0
	for idx, count := range answerCounts {
		if count >= answerCounts[highestIdx] {
			highestIdx = idx
		}
	}

	// Format winning response and account for ties
	winningAnswer := fmt.Sprintf("[%d] %s with %d votes", highestIdx, bot.pollAnswers[highestIdx], answerCounts[highestIdx])
	for idx, count := range answerCounts {
		if count == answerCounts[highestIdx] {
			winningAnswer += " | "
			winningAnswer += fmt.Sprintf("[%d] %s with %d votes", idx, bot.pollAnswers[idx], answerCounts[idx])
		}
	}
	bot.SendMessage("Poll Winner(s): %s", winningAnswer)

	logger.Info("Closing poll")
	bot.pollRunning = false
	bot.pollQuestion = ""
	bot.pollAnswers = []string{}
}

// sendEvent sends a Bot event to Write-enabled clients
func (bot *Bot) sendEvent(evt Event) {
	bot.chatClient.Channel() <- evt
}

// Start listening for Events on the inbound channel and broadcast out
// to the Handlers
func (bot *Bot) listen() {
	// Spawn goroutines for Handlers
	for _, consumer := range bot.consumers {
		go consumer.Listen()
	}

	bot.listening = true

	for {
		select {
		case evt := <-bot.events:
			bot.Mutex.Lock()
			for _, consumer := range bot.consumers {
				consumer.Receive(evt)
			}
			bot.Mutex.Unlock()
		default:
			// TODO QUIT message
		}
	}
}
