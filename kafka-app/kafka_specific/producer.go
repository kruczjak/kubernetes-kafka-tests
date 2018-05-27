package kafka_specific

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"time"
	"strconv"
	"fmt"
	"strings"
)

func StartProducer(brokers *string, maxMessages int) {
	config := sarama.NewConfig()

	brokerList := strings.Split(*brokers, ",")

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// close producer after finishing
	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	// shutdown correctly on SIGINT and SIGKILL
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)

	// main loop with producing messages
	var enqueued, errors int
	doneCh := make(chan struct{})
	go func() {
		for {
			if enqueued > maxMessages {
				break
			}

			strTime := strconv.Itoa(int(time.Now().Unix()))
			msg := &sarama.ProducerMessage{
				Topic: "important",
				Key:   sarama.StringEncoder(strTime),
				Value: sarama.StringEncoder("Something Cool"),
			}
			select {
			case producer.Input() <- msg:
				enqueued++
				fmt.Println("Produce message")
			case err := <-producer.Errors():
				errors++
				fmt.Println("Failed to produce message:", err)
			case <-signals:
				doneCh <- struct{}{}
			}
		}
	}()

	// wait until something goes to doneCh channel
	<-doneCh
	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
}
