package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventHub struct {
	reouteService   *RouteService
	mongoClient     *mongo.Client
	chDriverMoved   chan *DriverMovedEvent
	freightWriter   *kafka.Writer
	simulatirWriter *kafka.Writer
}

func NewEventHub(
	reouteService *RouteService,
	mongoClient *mongo.Client,
	chDriverMoved chan *DriverMovedEvent,
	freightWriter *kafka.Writer,
	simulatirWriter *kafka.Writer,
) *EventHub {
	return &EventHub{
		reouteService:   reouteService,
		mongoClient:     mongoClient,
		chDriverMoved:   chDriverMoved,
		freightWriter:   freightWriter,
		simulatirWriter: simulatirWriter,
	}
}

func (eh *EventHub) HandlerEvent(msg []byte) error {
	var baseEvent struct {
		EventName string `json:"event"`
	}
	err := json.Unmarshal(msg, &baseEvent)
	if err != nil {
		return fmt.Errorf("")
	}

	switch baseEvent.EventName {
	case "RouteCreated":
		var routeCreatedEvent RouteCreatedEvent
		json.Unmarshal(msg, &routeCreatedEvent)
		eh.handlerRouteCreated(&routeCreatedEvent)
	case "DeliveryStarted":
		var deliveryStartedEvent DeliveryStartedEvent
		json.Unmarshal(msg, &deliveryStartedEvent)
		err := DeliveryStartedHandler(&deliveryStartedEvent, eh.reouteService, eh.chDriverMoved)
		if err != nil {
			return err
		}
	}
	return nil
}
func (eh *EventHub) handlerRouteCreated(event *RouteCreatedEvent) error {
	freightEvent, err := RouteCreatedHandler(event, eh.reouteService)
	if err != nil {
		return err
	}
	value, err := json.Marshal(freightEvent)
	if err != nil {
		return err
	}

	log.Printf("publishing freight event: %s\n", string(value))
	err = eh.freightWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(freightEvent.RouteId),
		Value: value,
	})
	if err != nil {
		return err
	}

	return nil
}

func (eh *EventHub) sendDirections() {
	for {
		select {
		case driverMovedEvent := <-eh.chDriverMoved:
			value, err := json.Marshal(driverMovedEvent)
			if err != nil {
				return
			}
			log.Printf("publishing driver moved event: %s\n", string(value))
			err = eh.simulatirWriter.WriteMessages(context.Background(), kafka.Message{
				Key:   []byte(driverMovedEvent.RouteId),
				Value: value,
			})
			if err != nil {
				return
			}
		case <-time.After(500 * time.Millisecond):
			return
		}
	}
}

func (eh *EventHub) handlerDeliveryStarted(event *DeliveryStartedEvent) error {
	err := DeliveryStartedHandler(event, eh.reouteService, eh.chDriverMoved)
	if err != nil {
		return err
	}

	go eh.sendDirections()

	return nil
}
