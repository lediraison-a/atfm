package main

import (
	"atfm/app"
	"atfm/app/config"
	"atfm/app/models"
)

func main() {
	// _, err := server.NewServer()
	// if err != nil {
	//  	log.Print(err)
	// }
	// defer s.Inbound.Close()
	config := config.NewConfigDefault()
	// client, err := app.NewClient()
	// if err != nil {
	//  	panic(err)
	// }
	// instances := app.NewInstancePool(client.RpcClient)
	instances := app.NewInstancePool(nil)
	tui := app.NewTui(instances, config)
	err := tui.NewInstance(config.Start.StartDir, config.Start.StartBasepath, models.LOCALFM, true)
	if err != nil {
		panic(err)
	}
	tui.StartApp()
}
