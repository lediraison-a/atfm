package main

import (
	"atfm/app"
	"atfm/app/config"
	"atfm/app/models"
	"atfm/app/server"
	"log"
)

func main() {
	_, err := server.NewServer()
	if err != nil {
		log.Print(err)
	}
	// defer s.Inbound.Close()
	config := config.NewConfigDefault()
	client, err := app.NewClient()
	if err != nil {
		panic(err)
	}
	instances := app.NewInstancePool(client.RpcClient)
	tui := app.NewTui(instances, config)
	tui.NewInstance(config.Start.StartDir, config.Start.StartBasepath, models.LOCALFM, true)
	tui.StartApp()
}
