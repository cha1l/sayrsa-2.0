package main

import (
	"log"
	"os"

	sayrsa20 "github.com/cha1l/sayrsa-2.0"
	"github.com/cha1l/sayrsa-2.0/pkg/handler"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
	"github.com/cha1l/sayrsa-2.0/pkg/service"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	if err := InitConfig(); err != nil {
		log.Fatal("error while reading config")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("error while loading .env file")
	}

	c := repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBname:   viper.GetString("db.dbname"),
		Sslmode:  viper.GetString("db.sslmode"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	db, err := repository.NewDB(c)
	if err != nil {
		log.Fatal("failed to launch db...")
	}

	repository := repository.New(db)
	service := service.New(repository)
	handler := handler.New(service)

	apiserver := new(sayrsa20.APIServer)
	if err := apiserver.StartServer(viper.GetString("port"), handler.InitRoutes()); err != nil {
		log.Fatal(err)
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
