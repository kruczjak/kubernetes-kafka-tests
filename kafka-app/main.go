package main

import (
	"log"
	"flag"
	"./kafka_specific"
)

var (
	workerType = flag.String("worker_type", "producer", "Producer or consumer")
)
func main() {
	log.Printf("Starting app")
	flag.Parse()

	if *workerType == "producer" {
		kafka_specific.StartProducer()
	} else {
		kafka_specific.StartConsumer()
	}
}
