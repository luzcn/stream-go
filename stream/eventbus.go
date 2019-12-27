package stream

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

// Config records the fields that are used to identify the EventBus-Sub client
type Config struct {
	Endpoint  string
	AuthToken string
	Stream    string
	Client    string
	Version   string
}

type EventHandler interface {
	Handler(Message) error
}

// EventHandlerFunc is an adapter type to allow the use of ordinary functions
// as an EventHandler
type EventHandlerFunc func(Message) error

// Handle implements EventHandler for the EventHandlerFunc adapter type.
func (e EventHandlerFunc) Handle(m Message) error {
	return e(m)
}

type Message struct {
	Offset    int64           `json:"offset"`
	Partition int32           `json:"partition"`
	Body      json.RawMessage `json:"body"`
}

// An EventBus object is the client for connecting EventBus
type EventBus struct {
	config         Config
	state          string
	eventHandler   EventHandler
	startingOffset int64
	socket         *websocket.Conn
	store          offsetStore
}

func NewEventBus(config Config, eventHandler EventHandler) *EventBus {
	return &EventBus{
		config: config,
		// store:        store,
		eventHandler: eventHandler,
	}
}

func (eb *EventBus) connect() error {
	c, _, err := websocket.DefaultDialer.Dial(eb.config.Endpoint, nil)
	if err != nil {
		return err
	}
	eb.socket = c
	return nil
}

// Run starts the eventbus loop.
// When Run is called, the registered EventHandler will be called for each
// message in the stream.
// It returns a chan that the caller can wait on to receive errors during event
// streaming.
func (eb *EventBus) Run() chan error {

	done := make(chan error)
	go func() {
		defer close(done)
		defer func() {
			if eb.socket != nil {
				_ = eb.socket.Close()
			}
		}()

		for {
			if eb.socket == nil {
				err := eb.connect()
				if err != nil {
					done <- err
					return
				}
			}

			_, msg, err := eb.socket.ReadMessage()
			if err != nil {
				log.Println(err)
				_ = eb.socket.Close()
				eb.socket = nil
				done <- err
				continue
			}

			fmt.Println(string(msg))
		}
	}()

	return done
}
