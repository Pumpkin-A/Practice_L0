package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/segmentio/kafka-go"
)

// the topic and broker address are initialized as constants
const (
	topic          = "cats3"
	broker1Address = "localhost:9092"
	broker2Address = "localhost:9093"
	broker3Address = "localhost:9094"
)

var wg sync.WaitGroup

func produce(ctx context.Context) {
	// initialize a counter
	i := 0

	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address, broker2Address, broker3Address},
		Topic:   topic,
	})

	for range 10 {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),
			// create an arbitrary message payload for the value
			Value: []byte("this is message" + strconv.Itoa(i)),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// log a confirmation once the message is written
		fmt.Println("writes:", i)
		i++
		// sleep for a second
		time.Sleep(time.Second)
	}
}

// var isProducerEnded bool = false

func consume(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	ctx, cancel := context.WithCancel(ctx)
	var counter atomic.Int32
	go func() {
		defer cancel()
		for {
			time.Sleep(time.Second)
			if counter.Load() == 10 {
				return
			}
		}
	}()
	for i := range 5 {
		wg.Add(1)
		go func(i int) {
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{broker1Address, broker2Address, broker3Address},
				Topic:   topic,
				GroupID: "my-group",
			})
			for {
				// the `ReadMessage` method blocks until we receive the next event
				msg, err := r.ReadMessage(ctx)
				if err != nil {
					fmt.Println("could not read message " + err.Error())
					break
				}
				// after receiving the message, log its value
				fmt.Println("received: ", i, string(msg.Value))
				counter.Add(1)
			}
			fmt.Printf("reader %v was closed\n", i)
			wg.Done()
		}(i)
	}
}

func main() {
	// create a new context
	ctx := context.Background()
	// produce messages in a new go routine, since
	// both the produce and consume functions are
	// blocking
	consume(ctx)
	produce(ctx)
	wg.Wait()
	// time.Sleep(1*time.Minute)
}
