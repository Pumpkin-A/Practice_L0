package consumer

import (
	"context"
	"log/slog"
	"practiceL0_go_mod/config"
	"practiceL0_go_mod/internal/models"
	"sync"

	"github.com/segmentio/kafka-go"
)

type OrderManager interface {
	SaveOrder(msg []byte) error
}

type KafkaConsumer struct {
	OrderManager      OrderManager
	Topic             string
	Broker1Address    string
	Broker2Address    string
	Broker3Address    string
	NumberOfConsumers int
}

func New(cfg config.Config, om OrderManager) *KafkaConsumer {
	consumer := &KafkaConsumer{
		OrderManager:      om,
		Topic:             cfg.Kafka.Topic,
		Broker1Address:    cfg.Kafka.Broker1Address,
		Broker2Address:    cfg.Kafka.Broker2Address,
		Broker3Address:    cfg.Kafka.Broker3Address,
		NumberOfConsumers: cfg.Kafka.NumberOfConsumers,
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
			slog.Info("kafka consumer start", "consumerNumber", i)
			defer slog.Info("reader was closed", "consumerNumber", i)
			defer r.Close()

			for {
				// the `FetchMessage` method blocks until we receive the next event
				msg, err := r.FetchMessage(mainCtx)
				if err != nil {
					slog.Error("could not fetch message", "func", "Consumer: Run", "err", err.Error())
					break
				}

				err = c.OrderManager.SaveOrder(msg.Value)
				if err != nil {
					slog.Error("error with msg processing in consumer", "func", "Consumer: Run", "err", err.Error())
					if err != models.ErrorValidation {
						break
					}
				}

				err = r.CommitMessages(context.Background(), msg)
				if err != nil {
					slog.Error("error with kafka committing msg", "func", "Consumer: Run", "err", err.Error())
					break
				}
			}
		}(i)
	}
	wg.Wait()

}
