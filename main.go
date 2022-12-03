package main

import (
	"atfm/app"
	"atfm/app/config"
	"atfm/app/models"
	"atfm/app/server"
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
	instances := app.NewInstancePool(nil, config)
	tui := app.NewTui(instances, config)

	or := func(path string, content []models.FileInfo, selfDelete bool) error {
		tui.RefreshInstances(path, content, selfDelete)
		return nil
	}
	app.FileManagerService = server.NewFileManager()
	app.NotifyManagerService = server.NewNotifyManager(app.FileManagerService, or)

	err := tui.NewInstance(config.StartDir, config.StartBasepath, models.LOCALFM, true)
	if err != nil {
		panic(err)
	}
	tui.StartApp()
}
