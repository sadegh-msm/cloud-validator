package model

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws/session"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
)

type Resources struct {
	MongoDB          *mongo.Client
	MongoColl        *mongo.Collection
	S3Sess           *session.Session
	RabbitConnection *amqp.Connection
}

var Res Resources

type User struct {
	Name       string `json:"name" bson:"name"`
	Email      string `json:"email" bson:"email"`
	NationalId string `json:"nationalId" bson:"_id"`
	IP         string `json:"ip" bson:"ip"`
	Image1     string `json:"image1" bson:"image1"`
	Image2     string `json:"image2" bson:"image2"`
	State      string `json:"state" bson:"state"`
}

func ConnectMongo() (err error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://msmohamadi1380:13sadegh81@hw1-cloud.9hbuqq3.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

	Res.MongoDB, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Warnln(err)
		return err
	}

	Res.MongoColl = Res.MongoDB.Database("validator").Collection("users")
	return nil
}

func PingDB(client *mongo.Client) (err error) {
	if err = client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Warnln("error", err)
		return err
	}

	log.Infoln("database is fine!")
	return nil
}

func CloseConn(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Warnln(err)
	}
}

func Insert(name, email, nationalId, ip, image1, image2 string) error {
	err := Res.MongoColl.FindOne(context.TODO(), bson.D{{"_id", nationalId}}).Err()
	if err != mongo.ErrNoDocuments {
		log.Infoln("user is already registered: ", nationalId)
		return err
	}

	nationalId = base64.StdEncoding.EncodeToString([]byte(nationalId))

	doc := &User{
		Name:       name,
		Email:      email,
		NationalId: nationalId,
		IP:         ip,
		Image1:     image1,
		Image2:     image2,
		State:      "pending",
	}
	_, err = Res.MongoColl.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Panicln(err)
	}

	log.Infoln("user is created: ", nationalId, " name:", name)

	return nil
}

func Update(nationalId, state string) bool {
	update := bson.D{
		{"$set", bson.D{
			{"state", state},
		}},
	}
	_, err := Res.MongoColl.UpdateOne(context.TODO(), bson.D{{"_id", nationalId}}, update)
	if err != nil {
		log.Warnln("cant update users object")
		return false
	}

	return true
}

func GetAll() []User {
	cur, err := Res.MongoColl.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Warnln("cant find all users")
	}

	res := make([]User, 0)
	var doc User
	for cur.Next(context.TODO()) {
		err := cur.Decode(&doc)
		if err != nil {
			log.Panicln(err)
		}
		res = append(res, doc)
	}

	return res
}

func Find(nationalId string) *User {
	var doc User

	err := Res.MongoColl.FindOne(context.TODO(), bson.D{{"_id", nationalId}}).Decode(&doc)
	if err != nil {
		log.Warnln("user not found")
	}

	return &doc
}
