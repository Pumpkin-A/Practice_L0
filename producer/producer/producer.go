package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"practice_L0_producer_gomod/orders"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	topic          = "orders"
	broker1Address = "localhost:9092"
	// broker2Address = "localhost:9093"
	// broker3Address = "localhost:9094"
)

func Produce(ctx context.Context) {
	// initialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address},
		Topic:   topic,
	})

	for {
		order := orders.GenerateOrder()
		jsonOrder, err := json.Marshal(order)
		if err != nil {
			fmt.Printf("Error: %s", err)
		}

		err = w.WriteMessages(ctx, kafka.Message{
			Key:   []byte(order.OrderUID.String()),
			Value: []byte(jsonOrder),
		})
		if err != nil {
			fmt.Println("could not write message " + err.Error())
			return
		}

		// log a confirmation once the message is written
		fmt.Println("writes:", order)
		// sleep for a second
		time.Sleep(time.Second)
	}
}
