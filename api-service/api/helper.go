package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"hw1/api-service/configs"

	"context"
	"fmt"
	"hw1/api-service/model"
	"mime/multipart"
)

func ConnectS3() (err error) {
	model.Res.S3Sess, err = session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(configs.Conf.S3AccessKey, configs.Conf.S3SecretKey, ""),
		Region:      aws.String(configs.Conf.S3Region),
		Endpoint:    aws.String(configs.Conf.S3Endpoint),
	})
	if err != nil {
		log.Warnln("can not connect to s3", err)
		return err
	}

	log.Infoln("connected to S3 instance")
	return nil
}

func UploadS3(sess *session.Session, fileHeader *multipart.FileHeader, bucket string, ID string) string {
	uploader := s3manager.NewUploader(sess)
	file, err := fileHeader.Open()
	if err != nil {
		log.Warnln("cant open file")
		return ""
	}

	fullFileName := ID + fileHeader.Filename
	key := fmt.Sprintf("%s", fullFileName)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(configs.Conf.S3Bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Warnf("Unable to upload %q to %q, %v", fullFileName, bucket, err)
		return ""
	}
	log.Infof("Successfully uploaded %q to %q\n", fullFileName, bucket)

	return key
}

func listMyBuckets(sess *session.Session) {
	svc := s3.New(sess, &aws.Config{
		Region:   aws.String(configs.Conf.S3Region),
		Endpoint: aws.String(configs.Conf.S3Endpoint),
	})

	result, err := svc.ListBuckets(nil)

	if err != nil {
		log.Warnf("Unable to list buckets, %v", err)
	}

	log.Infoln("My buckets now are:")

	for _, b := range result.Buckets {
		log.Infoln(aws.StringValue(b.Name) + "\n")
	}
}

func ConnectMQ() (err error) {
	model.Res.RabbitConnection, err = amqp.Dial(configs.Conf.AmqpAddress)
	if err != nil {
		return err
	}
	return err
}

func WriteMQ(connection *amqp.Connection, message string) error {
	channel, _ := connection.Channel()

	msg := amqp.Publishing{
		DeliveryMode: 1,
		ContentType:  "text/plain",
		Body:         []byte(message),
	}

	err := channel.PublishWithContext(context.TODO(), "amq.topic", "ping", false, false, msg)
	if err != nil {
		log.Warnln("cant publish message to queue")
		return err
	}

	return nil
}

func CloseMQ(connection *amqp.Connection) {
	err := connection.Close()
	if err != nil {
		return
	}
}
