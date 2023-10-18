package configs

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"os"
)

type Config struct {
	S3AccessKey string
	S3SecretKey string
	S3Bucket string
	S3Endpoint string
	S3Region string
	MailGunApiKey string
	MailGunDomain string
	AmqpAddress string
	MongoAddress string
}

var Conf Config

func SetConf() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Conf.S3AccessKey = os.Getenv("S3_ACCESS_KEY")
	Conf.S3SecretKey = os.Getenv("S3_SECRET_KEY")
	Conf.S3Bucket = os.Getenv("S3_BUCKET")
	Conf.S3Endpoint = os.Getenv("S3_ENDPOINT")
	Conf.S3Region = os.Getenv("S3_REGION")
	Conf.MailGunApiKey = os.Getenv("MAILGUN_API_KEY")
	Conf.MailGunDomain = os.Getenv("MAILGUN_DOMAIN")
	Conf.AmqpAddress = os.Getenv("AMQP_ADDRESS")
	Conf.MongoAddress = os.Getenv("MONGO_ADDRESS")
}
