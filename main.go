package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/luzcn/stream/stream"
	"os"
)

type eventHandler struct {
}

func (e eventHandler) Handler(m stream.Message) error {

	fmt.Println(m.Body)
	return nil
}

func main() {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	eventbusConfig := stream.Config{
		Endpoint:  os.Getenv("EVENTBUS_URL"),
		AuthToken: os.Getenv("EVENTBUS_TOKEN"),
		Stream:    "api-events",
		Client:    "stream-go-test",
		Version:   "0.0.1",
	}

	eb := stream.NewEventBus(eventbusConfig, eventHandler{})

	<-eb.Run()

}
