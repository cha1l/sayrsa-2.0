package main

import (
	sayrsa20 "github.com/cha1l/sayrsa-2.0"
	"github.com/cha1l/sayrsa-2.0/pkg/handler"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
	"github.com/cha1l/sayrsa-2.0/pkg/service"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	if err := InitConfig(); err != nil {
		log.Fatal("error reading config")
	}

	var c repository.Config
	var ip string
	var port string

	if viper.GetString("start") == "localhost" {
		c = repository.Config{
			Host:     viper.GetString("db.host"),
			User:     viper.GetString("db.user"),
			Password: viper.GetString("db.password"),
			DBname:   viper.GetString("db.name"),
		}
		ip = "localhost"
		port = viper.GetString("port")
	} else {
		c = repository.Config{
			Host:     os.Getenv("DBHOST"),
			User:     os.Getenv("DBUSER"),
			Password: os.Getenv("DBPASS"),
			DBname:   os.Getenv("DBNAME"),
		}
		ip = os.Getenv("APP_IP")
		port = os.Getenv("APP_PORT")
	}

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

func InitConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
