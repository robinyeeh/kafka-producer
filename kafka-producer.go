package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage:  ./kafka_producer 127.0.0.1:9092 test_topic")
		return
	}

	address := []string{args[1]}
	topic := args[2]

	produceMessage(address, topic)
}

func produceMessage(address []string, topic string) {
	fmt.Printf("Produce kafka message, address : %s, topic : %s \n", address, topic)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		fmt.Printf("Create producer error : %s\n", err.Error())
		return
	}

	defer producer.AsyncClose()

	//send message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("producer_key"),
	}

	value := "This is a test message"
	for {
		fmt.Scanln(&value)
		msg.Value = sarama.ByteEncoder(value)
		fmt.Printf("Input [%s]\n", value)

		//send to chain
		producer.Input() <- msg

		select {
		case suc := <-producer.Successes():
			{
				fmt.Printf("Offset : %d, timestamp : %s", suc.Offset, suc.Timestamp.String())
			}
		case fail := <-producer.Errors():
			{
				fmt.Printf("Errors : %s\n", fail.Err.Error())
			}

		}

	}
}
