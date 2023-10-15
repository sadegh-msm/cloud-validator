package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"fmt"
	"mime/multipart"
	"time"
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

func UploadS3(sess *session.Session, fileHeader *multipart.FileHeader, bucket string, ID string) string {
	uploader := s3manager.NewUploader(sess)
	file, err := fileHeader.Open()
	key := fmt.Sprintf("%s", fileHeader.Filename+ID)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Warnln("Unable to upload %q to %q, %v", fileHeader.Filename, bucket, err)
	}
	log.Infoln("Successfully uploaded %q to %q\n", fileHeader.Filename, bucket)

	return key
}

func ConnectMQ() *amqp.Connection {
	url := "amqps://xgyeesmr:T-UTG1qOjoipEH5wB5xFoLPInQ7MpjYJ@sparrow.rmq.cloudamqp.com/xgyeesmr"
	connection, _ := amqp.Dial(url)

	return connection
}

func WriteMQ(connection *amqp.Connection, message string) error {
	timer := time.NewTicker(1 * time.Second)
	channel, _ := connection.Channel()

	for t := range timer.C {
		msg := amqp.Publishing{
			DeliveryMode: 1,
			Timestamp:    t,
			ContentType:  "text/plain",
			Body:         []byte(message),
		}
		mandatory, immediate := false, false
		err := channel.Publish("amq.topic", "ping", mandatory, immediate, msg)
		if err != nil {
			log.Warnln("cant publish message to queue")
			return err
		}
	}
	return nil
}
