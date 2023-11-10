package api

import (
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
	"hw1/validation-service/configs"

	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func ConnectS3() (err error) {
	Res.S3Sess, err = session.NewSession(&aws.Config{
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

func DownloadS3(sess *session.Session, bucket string, key string) *os.File {
	getwd, _ := os.Getwd()

	file, err := os.Create(getwd + "/validation-service/images/" + key)
	if err != nil {
		log.Warnf("Unable to open file %q, %v", key, err)
	}

	defer file.Close()

	s3Client := s3.New(sess)

	obj, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(configs.Conf.S3Bucket),
		Key:    aws.String(key),
	})

	_, err = io.Copy(file, obj.Body)
	if err != nil {
		log.Warnln("cant copy file")
	}

	if err != nil {
		log.Warnf("Unable to download item %q, %v", key, err)
	}

	return file
}

func ConnectMQ() (err error) {
	Res.RabbitConnection, err = amqp.Dial(configs.Conf.AmqpAddress)
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
	opts := options.Client().ApplyURI(configs.Conf.MongoAddress).SetServerAPIOptions(serverAPI)

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

func faceDetection(file *os.File) (string, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	file, err := os.Open(file.Name())
	if err != nil {
		return "", err
	}

	defer file.Close()

	part, err := writer.CreateFormFile("image", file.Name())
	if err != nil {
		log.Warnln("Error creating form file:", err)
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Warnln("Error copying image to request:", err)
		return "", err
	}
	writer.Close()

	url := "https://api.imagga.com/v2/faces/detections?return_face_id=1"
	request, err := http.NewRequest("POST", url, bytes.NewReader(requestBody.Bytes()))
	if err != nil {
		log.Warnln("Error creating request:", err)
		return "", err
	}

	request.SetBasicAuth(configs.Conf.ImaggaApiKey, configs.Conf.ImaggaSecretKey)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Warnln("Error making request:", err)
		return "", err
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, response.Body)
	if err != nil {
		log.Warnln("Error reading response:", err)
		return "", err
	}

	log.Infoln("Response:", buf.String())

	id := ParseFaceIdJSON(buf.String())
	return id, nil
}

func FaceSimilarity(face1, face2 string) int {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api.imagga.com/v2/faces/similarity?face_id="+face1+"&second_face_id="+face2, nil)
	req.SetBasicAuth(configs.Conf.ImaggaApiKey, configs.Conf.ImaggaSecretKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error when sending request to the server")
		return 0
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	res := ParseScoreJSON(string(respBody))

	return res
}

func Update(nationalId, state string) bool {
	update := bson.D{
		{"$set", bson.D{
			{"state", state},
		}},
	}

	nationalId = base64.StdEncoding.EncodeToString([]byte(nationalId))
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
