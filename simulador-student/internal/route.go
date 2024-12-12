package internal

import (
	"context"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Directions struct {
	Lat float64 `bson:"lat"`
	Lng float64 `bson:"lng"`
}

type Route struct {
	Id           string       `bson:"_id" json:"id"`
	Distance     int          `bson:"distance" json:"distance"`
	Directions   []Directions `bson:"diretions" json:"diretions"`
	FreightPrice float64      `bson:"freight_price" json:"freight_price"`
}

func NewRoute(id string, distance int, directions []Directions) Route {
	return Route{
		Id:         id,
		Distance:   distance,
		Directions: directions,
	}
}

type RouteService struct {
	mongo          *mongo.Client
	freightService *FreightService
}

func NewRouteService(conn *mongo.Client, freightService *FreightService) *RouteService {
	return &RouteService{
		mongo:          conn,
		freightService: freightService,
	}
}

type FreightService struct{}

func (fs *FreightService) Calculate(distance int) float64 {
	return math.Floor((float64(distance)*0.15+0.3)*100) / 100
}
func NewFreightService() *FreightService {
	return &FreightService{}
}

func (rs *RouteService) CreateRoute(route Route) (Route, error) {

	route.FreightPrice = rs.freightService.Calculate(route.Distance)

	update := bson.M{
		"$set": bson.M{
			"distance":      route.Distance,
			"direction":     route.Directions,
			"freigth_price": route.FreightPrice,
		},
	}
	filter := bson.M{"_id": route.Id}

	options := options.Update().SetUpsert(true)

	_, err := rs.mongo.Database("golang").Collection("routes").UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		return Route{}, err
	}

	return route, err
}

func (rs *RouteService) GetRoute(id string) (Route, error) {
	var route Route

	filter := bson.M{"_id": id}

	err := rs.mongo.Database("golang").Collection("routes").FindOne(context.Background(), filter).Decode(&route)

	return route, err
}
