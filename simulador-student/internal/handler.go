package internal

import "time"

type RouteCreatedEvent struct {
	EventName  string       `json:"event"`
	RouteId    string       `json:"id"`
	Distance   int          `json:"distance"`
	Directions []Directions `json:"directions"`
}

func NewRouteCreatedEvent(routeID string, distance int, directions []Directions) *RouteCreatedEvent {
	return &RouteCreatedEvent{
		EventName:  "RouteCreated",
		RouteId:    routeID,
		Distance:   distance,
		Directions: directions,
	}
}

type FreightCalculatedEvent struct {
	EventName string  `json:"event"`
	RouteId   string  `json:"route_id"`
	Amount    float64 `json:"amount"`
}

func NewFreightCalculatedEvent(routeID string, amount float64) *FreightCalculatedEvent {
	return &FreightCalculatedEvent{
		EventName: "FreightCalculated",
		RouteId:   routeID,
		Amount:    amount,
	}
}

type DeliveryStartedEvent struct {
	EventName string `json:"event"`
	RouteId   string `json:"route_id"`
}

func NewDeliveryStartedEvent(routeId string) *DeliveryStartedEvent {
	return &DeliveryStartedEvent{
		EventName: "DeliveryStarted",
		RouteId:   routeId,
	}
}

type DriverMovedEvent struct {
	EventName string  `json:"event"`
	RouteId   string  `json:"route_id"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
}

func NewDriverMovedEvent(routeId string, lat, lng float64) *DriverMovedEvent {
	return &DriverMovedEvent{
		EventName: "DriverMoved",
		RouteId:   routeId,
		Lat:       lat,
		Lng:       lng,
	}
}

func RouteCreatedHandler(event *RouteCreatedEvent, routeService *RouteService) (*FreightCalculatedEvent, error) {
	route := NewRoute(event.RouteId, event.Distance, event.Directions)
	createdRoute, err := routeService.CreateRoute(route)
	if err != nil {
		return nil, err
	}
	freightEvent := NewFreightCalculatedEvent(createdRoute.Id, createdRoute.FreightPrice)
	return freightEvent, nil
}

func DeliveryStartedHandler(event *DeliveryStartedEvent, routeService *RouteService, ch chan<- *DriverMovedEvent) error {
	route, err := routeService.GetRoute(event.RouteId)
	if err != nil {
		return err
	}
	driverMovedEvent := NewDriverMovedEvent(event.RouteId, 0, 0)
	for _, direction := range route.Directions {
		driverMovedEvent.Lat = direction.Lat
		driverMovedEvent.Lng = direction.Lng
		time.Sleep(time.Second)
		ch <- driverMovedEvent
	}
	return nil
}
