package Employees

import (

	"context"
	"log"
	"mongox"
	"redisx"
)



type UserControl struct {
	MongoClient *mongox.Client
	Redis       *redisx.Client

}

func NewUserControl() *UserControl {

	//connect to mongodb
	ctx := context.Background()
	Mongoclient, err := mongox.Connect(ctx, Configuration.Mongo)
	if err != nil {
		log.Fatal("Mongo connection failed:", err)
	}
	if err := Mongoclient.Ping(ctx); err != nil {
		log.Fatal("Mongo not reachable:", err)
	}
	//connect to redis
	redisclient, err := redisx.New(Configuration.Redis)
	if err != nil {
		log.Fatal("Error connecting to redis: ", err)
	}
	UC := &UserControl{

		MongoClient: Mongoclient,
		Redis:       redisclient,
	}
	return UC
}
