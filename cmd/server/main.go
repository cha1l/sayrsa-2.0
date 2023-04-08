package main

import (
	sayrsa20 "github.com/cha1l/sayrsa-2.0"
	"github.com/cha1l/sayrsa-2.0/pkg/handler"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
	"github.com/cha1l/sayrsa-2.0/pkg/service"
	"log"
	"os"
)

func main() {
	var c repository.Config
	var ip string
	var port string

	c = repository.Config{
		Host:     os.Getenv("DBHOST"),
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASS"),
		DBname:   os.Getenv("DBNAME"),
	}
	ip = os.Getenv("APP_IP")
	port = os.Getenv("APP_PORT")

	db, err := repository.NewDB(c)
	if err != nil {
		log.Fatal("failed to launch db...")
	}

	repo := repository.New(db)
	serv := service.New(repo)
	hand := handler.New(serv)

	apiserver := new(sayrsa20.APIServer)
	if err := apiserver.StartServer(hand.InitRoutes(), ip, port); err != nil {
		log.Fatal(err)
	}
}
