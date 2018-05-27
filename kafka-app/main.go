package main

import (
	"log"
	"flag"
	"./kafka_specific"
)

var (
	workerType = flag.String("worker_type", "producer", "Producer or consumer")
	brokers = flag.String("brokers", "", "Brokers addresses")
	producerMessages = flag.Int("messages", 1000000, "Number of producer messages")
)
func main() {
	log.Printf("Starting app")
	flag.Parse()

	if *workerType == "producer" {
		kafka_specific.StartProducer(brokers, *producerMessages)
	} else {
		kafka_specific.StartConsumer(brokers)
	}
}
