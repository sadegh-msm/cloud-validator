package main

import "hw1/api-service/api"

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

	e.Logger.Fatal(e.Start(s.Host + s.Port))
	//api.ConnectS3()
}
