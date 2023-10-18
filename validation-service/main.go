package main

import (
	log "github.com/sirupsen/logrus"

	"hw1/validation-service/api"
)

func main() {
	err := api.ConnectMongo()
	err = api.ConnectS3()
	err = api.ConnectMQ()
	api.InitChannel()
	if err != nil {
		log.Warnln(err)
		panic(err)
	}

	defer api.CloseConn(api.Res.MongoDB)
	defer api.CloseMQ(api.Res.RabbitConnection)

	//api.FaceSimilarity()

	for {
		api.Validate()
	}
}
