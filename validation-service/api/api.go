package api

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

func ConnectS3() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("ee82ad29-9bec-40f7-a64c-d854390c51a2", "24e63f33c5255a3862f0e7f83d6d37d519e2c55489b855f5dbba7bc5b41a45c4", ""),
		Region:      aws.String("default"),
		Endpoint:    aws.String("https://hw1-pic.s3.ir-thr-at1.arvanstorage.ir"),
	})

	if err != nil {
		log.Warnln("can not connect to s3", err)
	}
	log.Infoln("connected to S3 instance")

	return sess
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

func ConnectMQ() *amqp.Connection {
	url := "amqps://xgyeesmr:T-UTG1qOjoipEH5wB5xFoLPInQ7MpjYJ@sparrow.rmq.cloudamqp.com/xgyeesmr"
	connection, _ := amqp.Dial(url)

	return connection
}

func ReadMQ(connection *amqp.Connection) (string, error) {
	channel, _ := connection.Channel()

	q, _ := channel.QueueDeclare("user_ids", false, true, false, true, nil)
	err := channel.QueueBind(q.Name, "#", "amq.topic", false, nil)
	if err != nil {
		log.Warnln("cant bind to queue")
		return "", err
	}

	messages, _ := channel.Consume(q.Name, "", false, false, false, true, nil)
	for msg := range messages {
		fmt.Println("Body:", string(msg.Body), "Timestamp:", msg.Timestamp)
		msg.Ack(false)

		return string(msg.Body), nil
	}
	return "", nil
}
