package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mailgun/mailgun-go"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	apiKey    = "acc_fccb030428279de"
	secretKey = "b1a6e4450232c144c55495563d5599c0"
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

	s3Client := s3.New(sess)

	obj, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	_, err = io.Copy(file, obj.Body)
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Warnf("Unable to download item %q, %v", key, err)
	}

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

func SendMail(stage, receiver, domain, apiKey string) {
	mg := mailgun.NewMailgun(domain, apiKey)
	sender := "msmohamadi1380@gmail.com"
	subject := "Validation result"
	body := fmt.Sprintf("Your request for validating your information is on stage %s. Contact the admins if ypu have troubles for siginig in.", stage)

	sendMessage(mg, sender, subject, body, receiver)
}

func sendMessage(mg mailgun.Mailgun, sender, subject, body, recipient string) {
	message := mg.NewMessage(sender, subject, body, recipient)

	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("send main, also ID: %s Resp: %s\n", id, resp)
}

func faceDetection(file *os.File) {
	var requestBody *bytes.Buffer
	writer := multipart.NewWriter(requestBody)
	// Add the image file to the request
	part, err := writer.CreateFormFile("image", filepath.Base(file.Name()))
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	log.Infoln("body", requestBody.String())

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying image to request:", err)
		return
	}

	// Close the writer to finalize the request body
	writer.Close()

	// Create the HTTP request
	url := "https://api.imagga.com/v2/faces/detections"
	request, err := http.NewRequest("POST", url, bytes.NewReader(requestBody.Bytes()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the necessary headers for authentication and content type
	request.SetBasicAuth(apiKey, secretKey)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	log.Infoln("body end", request.Body)
	// Make the POST request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer response.Body.Close()

	log.Infoln(response.Status)
	// Read and print the response
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	log.Infoln("Response:", buf.String())
}

func FaceSimilarity() {
	output, err := exec.Command("ls").Output()
	if err != nil {
		return
	}
	log.Infoln(output)
}
