package consumer

import (
	"context"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

const (
	topic          = "orders"
	broker1Address = "localhost:9092"
	broker2Address = "localhost:9093"
	broker3Address = "localhost:9094"
)

type KafkaConsumer struct {
	WG *sync.WaitGroup
}

func New() *KafkaConsumer {
	var wg sync.WaitGroup
	return &KafkaConsumer{
		WG: &wg,
	}
}

func (c *KafkaConsumer) Consume(ctx context.Context) {
	// ctx, cancel := context.WithCancel(ctx)
	// ch := make(chan struct{})
	// // defer close(ch)

	// go func() {
	// 	defer cancel()
	// 	<-ch
	// }()

	for i := range 5 {
		c.WG.Add(1)
		go func(i int) {
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{broker1Address, broker2Address, broker3Address},
				Topic:   topic,
				GroupID: "bank",
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
			}
			fmt.Printf("reader %v was closed\n", i)
			c.WG.Done()
		}(i)
	}
}
