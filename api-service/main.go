package main

import (
	log "github.com/sirupsen/logrus"

	"hw1/api-service/api"
	"hw1/api-service/model"
)

type Server struct {
	Host string
	Port string
}

func main() {
	s := Server{
		"localhost",
		":8080",
	}
	e := api.SetupRouter()

	err := model.ConnectMongo()
	err = api.ConnectS3()
	err = api.ConnectMQ()
	if err != nil {
		log.Warnln(err)
		panic(err)
	}

	defer model.CloseConn(model.Res.MongoDB)
	defer api.CloseMQ(model.Res.RabbitConnection)

	e.Logger.Fatal(e.Start(s.Host + s.Port))
}
