package consumer

import (
	"context"
	"fmt"
	"practiceL0_go_mod/config"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Topic             string
	Broker1Address    string
	Broker2Address    string
	Broker3Address    string
	NumberOfConsumers int
	wg                *sync.WaitGroup
	ordersChan        chan []byte
}

func New(cfg config.Config) *KafkaConsumer {
	var wg sync.WaitGroup
	ch := make(chan []byte, cfg.Kafka.NumberOfConsumers)
	return &KafkaConsumer{
		Topic:             cfg.Kafka.Topic,
		Broker1Address:    cfg.Kafka.Broker1Address,
		Broker2Address:    cfg.Kafka.Broker2Address,
		Broker3Address:    cfg.Kafka.Broker3Address,
		NumberOfConsumers: cfg.Kafka.NumberOfConsumers,
		wg:                &wg,
		ordersChan:        ch,
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

	for i := range c.NumberOfConsumers {
		c.wg.Add(1)
		go func(i int) {
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{c.Broker1Address, c.Broker2Address, c.Broker3Address},
				Topic:   c.Topic,
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
