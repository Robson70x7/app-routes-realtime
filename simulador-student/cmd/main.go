package main

import (
	"context"
	"fmt"
	"log"

	"github.com/robson70x7/app-routes-realtime/simulador-student/internal"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoUrl := "mongodb://root:root@db:27017/golang?authSource=admin&directConnection=true"

	mongoConnection, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		panic(err)
	}

	freightService := internal.NewFreightService()
	routeService := internal.NewRouteService(mongoConnection, freightService)

	chDriverMoved := make(chan *internal.DriverMovedEvent)
	kafkdaBroker := "localhost:9092"
	freightWriter := kafka.Writer{
		Addr:     kafka.TCP(kafkdaBroker),
		Topic:    "freight",
		Balancer: &kafka.LeastBytes{},
	}
	simulatirWriter := kafka.Writer{
		Addr:     kafka.TCP(kafkdaBroker),
		Topic:    "simulator",
		Balancer: &kafka.LeastBytes{},
	}

	routeReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkdaBroker},
		Topic:   "route",
		GroupID: "simulator",
	})

	hub := internal.NewEventHub(routeService, mongoConnection, chDriverMoved, &freightWriter, &simulatirWriter)

	for {
		msg, err := routeReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}

		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

		go func(msg []byte) {
			err = hub.HandlerEvent(msg)
			if err != nil {
				log.Printf("error reading message: %v", err)
			}
		}(msg.Value)
	}

}
