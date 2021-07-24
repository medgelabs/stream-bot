package bot

import (
	"errors"
	"fmt"
	"medgebot/cache"
	"medgebot/logger"
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
	pollAnswers  []PollAnswer
}

// PollAnswer keeps track of an Answer label and number of votes for that answer
type PollAnswer struct {
	Answer string
	Count  int
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

	bot.ClearPoll()
	logger.Info("Starting new Poll: %s", question)
	bot.pollRunning = true
	bot.pollQuestion = question

	pollAnswers := make([]PollAnswer, len(answers))
	for idx, answer := range answers {
		pollAnswers[idx] = PollAnswer{
			Answer: answer,
			Count:  0,
		}
	}
	bot.pollAnswers = pollAnswers

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
	for idx, pollAnswer := range bot.pollAnswers {
		formattedAnswers += fmt.Sprintf("%d: %s | ", idx+1, pollAnswer.Answer)
	}

	bot.SendMessage("Poll started! Type a number only in chat to vote! Question: %s - | %s", bot.pollQuestion, formattedAnswers)
}

// AddPollVote increments the Count for the given Answer key
func (bot *Bot) AddPollVote(key int) {
	if key < 1 || key > len(bot.pollAnswers) {
		// Invalid vote. Skip
		return
	}

	bot.pollAnswers[key-1].Count++
}

// GetPollState returns the current question, answers, and vote counts for each answer
func (bot *Bot) GetPollState() (question string, answers []PollAnswer) {
	return bot.pollQuestion, bot.pollAnswers
}

// closePoll ends an active poll
func (bot *Bot) closePoll() {
	highestIdx := -1
	winningAnswer := PollAnswer{}
	for idx, answer := range bot.pollAnswers {
		if answer.Count >= winningAnswer.Count {
			highestIdx = idx
			winningAnswer = answer
		}
	}

	// If no votes - exit immediately
	if highestIdx == -1 {
		bot.SendMessage("No poll winner")
		bot.ClearPoll()
		return
	}

	// Format winning response and account for ties
	winnerStr := fmt.Sprintf("[%d] %s with %d votes", highestIdx+1, winningAnswer.Answer, winningAnswer.Count)
	for idx, answer := range bot.pollAnswers {
		if idx == highestIdx {
			continue
		}

		if answer.Count == winningAnswer.Count {
			winnerStr += " | "
			winnerStr += fmt.Sprintf("[%d] %s with %d votes", idx+1, answer.Answer, answer.Count)
		}
	}
	bot.SendMessage("Poll Winner(s): %s", winnerStr)

	// Ensure poll clears
	logger.Info("Closing poll")
	bot.ClearPoll()
}

// ClearPoll clears out any existing Poll state in the Bot and the dataStore
func (bot *Bot) ClearPoll() {
	bot.pollRunning = false
	bot.pollQuestion = ""
	bot.pollAnswers = []PollAnswer{}
	bot.dataStore.Clear("voters")
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
