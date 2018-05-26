package kafka_specific

import (
	"github.com/Shopify/sarama"
	"os"
	"os/signal"
	"fmt"
	"log"
)

func StartConsumer() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokerList := []string{"localhost:9092"}

	master, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	// consumer on topic
	topic := "important"
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	// shutdown correctly on SIGINT and SIGKILL
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)

	msgCount := 0
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Println("Received messages", string(msg.Key), string(msg.Value))
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	// wait until something goes to doneCh channel
	<-doneCh
	log.Println("Processed", msgCount, "messages")
}
