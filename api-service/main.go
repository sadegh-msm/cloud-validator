package main

import (
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

	model.ConnectMongo()
	defer model.CloseConn(model.DB)

	e.Logger.Fatal(e.Start(s.Host + s.Port))
}
