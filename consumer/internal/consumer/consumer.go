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
	wg         *sync.WaitGroup
	ordersChan chan []byte
}

func New() *KafkaConsumer {
	var wg sync.WaitGroup
	ch := make(chan []byte, 5)
	return &KafkaConsumer{
		wg:         &wg,
		ordersChan: ch,
	}
}

func (c *KafkaConsumer) GetOrdersChan() chan []byte {
	return c.ordersChan
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	// ctx, cancel := context.WithCancel(ctx)
	// ch := make(chan struct{})
	// // defer close(ch)

	// go func() {
	// 	defer cancel()
	// 	<-ch
	// }()

	for i := range 5 {
		c.wg.Add(1)
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
				// fmt.Println("received: ", i, string(msg.Value))
				c.ordersChan <- msg.Value
			}
			fmt.Printf("reader %v was closed\n", i)
			c.wg.Done()
		}(i)
	}
}
