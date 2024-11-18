package main

import (
	"context"
	"practice_L0_producer_gomod/producer"
)

func main() {
	ctx := context.Background()
	producer.Produce(ctx)
}
