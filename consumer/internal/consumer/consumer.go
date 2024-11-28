package consumer

import (
	"context"
	"fmt"
	"log"
	"practiceL0_go_mod/config"
	"sync"

	"github.com/segmentio/kafka-go"
)

type TransactionManager interface {
	AddConsumedOrdersToDBAndCache(msg []byte) error
}

type KafkaConsumer struct {
	TransactionManager TransactionManager
	Topic              string
	Broker1Address     string
	Broker2Address     string
	Broker3Address     string
	NumberOfConsumers  int
	wg                 *sync.WaitGroup
}

func New(cfg config.Config, tm TransactionManager) *KafkaConsumer {
	var wg sync.WaitGroup
	consumer := &KafkaConsumer{
		TransactionManager: tm,
		Topic:              cfg.Kafka.Topic,
		Broker1Address:     cfg.Kafka.Broker1Address,
		Broker2Address:     cfg.Kafka.Broker2Address,
		Broker3Address:     cfg.Kafka.Broker3Address,
		NumberOfConsumers:  cfg.Kafka.NumberOfConsumers,
		wg:                 &wg,
	}
	consumer.Start(context.Background())

	return consumer
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	for i := range c.NumberOfConsumers {
		c.wg.Add(1)
		go func(i int) {
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{c.Broker1Address, c.Broker2Address, c.Broker3Address},
				Topic:   c.Topic,
				GroupID: "bank",
			})
			for {
				msg, err := r.FetchMessage(ctx)
				if err != nil {
					log.Println("could not fetch message " + err.Error())
					break
				}
				err = c.TransactionManager.AddConsumedOrdersToDBAndCache(msg.Value)
				if err != nil {
					log.Println("error with msg processing in consumer", err.Error())
				}

				err = r.CommitMessages(ctx, msg)
				if err != nil {
					log.Println("error with kafka committing msg", err.Error())
					break
				}
			}

			fmt.Printf("reader %v was closed\n", i)
			c.wg.Done()
		}(i)
	}
}
