package api

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Resources struct {
	MongoDB          *mongo.Client
	MongoColl        *mongo.Collection
	S3Sess           *session.Session
	RabbitConnection *amqp.Connection
	RabbitChannel    *amqp.Channel
	RabbitQueue      amqp.Queue
	RabbitDelivery   <-chan amqp.Delivery
}

type User struct {
	Name       string `json:"name" bson:"name"`
	Email      string `json:"email" bson:"email"`
	NationalId string `json:"nationalId" bson:"_id"`
	IP         string `json:"ip" bson:"ip"`
	Image1     string `json:"image1" bson:"image1"`
	Image2     string `json:"image2" bson:"image2"`
	State      string `json:"state" bson:"state"`
}

type FaceDetectionResponse struct {
	Result struct {
		Faces []struct {
			FaceID string `json:"face_id"`
		} `json:"faces"`
	} `json:"result"`
}

type ScoreResponse struct {
	Result struct {
		Score float64 `json:"score"`
	} `json:"result"`
}

var Res Resources

func InitChannel() {
	Res.RabbitChannel, _ = Res.RabbitConnection.Channel()
	Res.RabbitQueue, _ = Res.RabbitChannel.QueueDeclare("user_ids", false, true, false, true, nil)

	err := Res.RabbitChannel.QueueBind(Res.RabbitQueue.Name, "#", "amq.topic", false, nil)
	if err != nil {
		log.Warnln("cant bind to queue")
	}
	Res.RabbitDelivery, _ = Res.RabbitChannel.Consume(Res.RabbitQueue.Name, "", false, false, false, true, nil)
}

func ParseFaceIdJSON(jsonData string) string {
	var response FaceDetectionResponse

	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return ""
	}

	if len(response.Result.Faces) > 0 {
		faceID := response.Result.Faces[0].FaceID
		return faceID
	}

	return ""
}

func ParseScoreJSON(jsonData string) int {
	var response ScoreResponse

	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return 0
	}

	score := response.Result.Score

	return int(score)
}
