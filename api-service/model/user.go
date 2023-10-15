package model

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
)

type User struct {
	Name       string                `json:"name"`
	Email      string                `json:"email"`
	UserId     string                `json:"userId"`
	NationalId string                `json:"nationalId"`
	IP         string                `json:"IP"`
	Image1     string `json:"image1"`
	Image2     string `json:"image2"`
	State      string                `json:"state"`
}

var DB *mongo.Client

func ConnectMongo() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://msmohamadi1380:13sadegh81@hw1-cloud.9hbuqq3.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	var err error
	DB, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Warnln(err)
	}
}

func PingDB(client *mongo.Client) {
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Warnln(err)
	}
	log.Infoln("database is fine!")
}

func CloseConn(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Warnln(err)
	}
}

func Insert() {

}