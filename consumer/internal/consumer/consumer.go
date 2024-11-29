package consumer

import (
	"context"
	"fmt"
	"log"
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/models"
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
}

func New(cfg config.Config, tm TransactionManager) *KafkaConsumer {
	consumer := &KafkaConsumer{
		TransactionManager: tm,
		Topic:              cfg.Kafka.Topic,
		Broker1Address:     cfg.Kafka.Broker1Address,
		Broker2Address:     cfg.Kafka.Broker2Address,
		Broker3Address:     cfg.Kafka.Broker3Address,
		NumberOfConsumers:  cfg.Kafka.NumberOfConsumers,
	}

	return consumer
}

func (c *KafkaConsumer) Run(mainCtx context.Context) {
	var wg sync.WaitGroup

	for i := range c.NumberOfConsumers {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{c.Broker1Address, c.Broker2Address, c.Broker3Address},
				Topic:   c.Topic,
				GroupID: "bank",
			})
			defer fmt.Printf("reader %d was closed\n", i)
			defer r.Close()

			for {
				// the `FetchMessage` method blocks until we receive the next event
				msg, err := r.FetchMessage(mainCtx)
				if err != nil {
					log.Println("could not fetch message " + err.Error())
					break
				}

				err = c.TransactionManager.AddConsumedOrdersToDBAndCache(msg.Value)
				if err != nil {
					log.Println("error with msg processing in consumer" + err.Error())
					if err != models.ErrorValidation {
						break
					}
				}

				err = r.CommitMessages(context.Background(), msg)
				if err != nil {
					log.Println("error with kafka committing msg" + err.Error())
					break
				}
			}
		}(i)
	}
	wg.Wait()

}
