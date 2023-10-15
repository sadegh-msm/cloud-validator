package model

import (
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
)

type User struct {
	Name       string `json:"name" bson:"name"`
	Email      string `json:"email" bson:"email"`
	NationalId string `json:"nationalId" bson:"nationalId"`
	IP         string `json:"ip" bson:"ip"`
	Image1     string `json:"image1" bson:"image1"`
	Image2     string `json:"image2" bson:"image2"`
	State      string `json:"state" bson:"state"`
}

var (
	DB   *mongo.Client
	Coll *mongo.Collection
)

func ConnectMongo() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://msmohamadi1380:13sadegh81@hw1-cloud.9hbuqq3.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	var err error
	DB, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Warnln(err)
	}

	Coll = DB.Database("validator").Collection("users")
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

func Insert(name, email, nationalId, ip, image1, image2 string) bool {
	err := Coll.FindOne(context.TODO(), bson.D{{"nationalId", nationalId}}).Err()
	if err != mongo.ErrNoDocuments {
		log.Infoln("user is already registered: ", nationalId)
		return false
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
	_, err = Coll.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Panicln(err)
	}

	log.Infoln("user is created: ", nationalId, " name:", name)

	return true
}

func Update(nationalId, state string) bool {
	update := bson.D{
		{"$set", bson.D{
			{"state", state},
		}},
	}
	_, err := Coll.UpdateOne(context.TODO(), bson.D{{"nationalId", nationalId}}, update)
	if err != nil {
		log.Warnln("cant update users object")
		return false
	}

	return true
}

func GetAll() []User {
	cur, err := Coll.Find(context.TODO(), bson.D{})
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

	err := Coll.FindOne(context.TODO(), bson.D{{"nationalId", nationalId}}).Decode(&doc)
	if err != nil {
		log.Warnln("user not found")
	}

	return &doc
}