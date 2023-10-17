package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func ConnectS3() (err error) {
	Res.S3Sess, err = session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("ee82ad29-9bec-40f7-a64c-d854390c51a2", "24e63f33c5255a3862f0e7f83d6d37d519e2c55489b855f5dbba7bc5b41a45c4", ""),
		Region:      aws.String("default"),
		Endpoint:    aws.String("https://hw1-pic.s3.ir-thr-at1.arvanstorage.ir"),
	})
	if err != nil {
		log.Warnln("can not connect to s3", err)
		return err
	}

	log.Infoln("connected to S3 instance")
	return nil
}

func DownloadS3(sess *session.Session, bucket string, key string) *os.File {
	file, err := os.Create(key)
	if err != nil {
		log.Warnf("Unable to open file %q, %v", key, err)
	}

	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		log.Warnf("Unable to download item %q, %v", key, err)
	}

	log.Infoln("Downloaded", file.Name(), numBytes, "bytes")

	return file
}

func ConnectMQ() (err error) {
	url := "amqps://xgyeesmr:T-UTG1qOjoipEH5wB5xFoLPInQ7MpjYJ@sparrow.rmq.cloudamqp.com/xgyeesmr"
	Res.RabbitConnection, err = amqp.Dial(url)
	if err != nil {
		return err
	}
	return err
}

func ReadMQ() (string, error) {
	for msg := range Res.RabbitDelivery {
		fmt.Println("Body:", string(msg.Body), "Timestamp:", msg.Timestamp)
		msg.Ack(false)

		return string(msg.Body), nil
	}
	return "", nil
}

func CloseMQ(connection *amqp.Connection) {
	err := connection.Close()
	if err != nil {
		return
	}
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

func Find(nationalId string) *User {
	var doc User

	nationalId = base64.StdEncoding.EncodeToString([]byte(nationalId))
	err := Res.MongoColl.FindOne(context.TODO(), bson.D{{"_id", nationalId}}).Decode(&doc)
	if err != nil {
		log.Warnln("user not found")
	}

	return &doc
}
