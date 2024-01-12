package main

import (
	"hotel-project/server"
	"hotel-project/util"
	"log"
)

func main() {
	envConfig, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	server := server.NewServer(&envConfig)
	app := server.RunServer()
	app.Listen(envConfig.HTTPServerAddress)
}
