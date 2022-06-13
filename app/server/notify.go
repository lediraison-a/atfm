package server

import (
	"atfm/app/models"
	"log"

	"github.com/rjeczalik/notify"
)

type NotifyManager struct {
	listeners map[int]chan notify.EventInfo

	filemanager *FileManager
}

func NewNotifyManager() *NotifyManager {
	return &NotifyManager{
		listeners: map[int]chan notify.EventInfo{},
	}
}

func (s *NotifyManager) SuscribeRefresh(instanceId int, path string, refreshFunc func([]models.FileInfo)) {
	pc, ok := s.listeners[instanceId]
	if ok {
		notify.Stop(pc)
		delete(s.listeners, instanceId)
	}
	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch(path, c, notify.Create, notify.Remove); err != nil {
		log.Fatal(err)
	}
	s.listeners[instanceId] = c

	go func() {
		defer notify.Stop(c)
		for {
			ei := <-c
			log.Println("Got event:", ei)
			var dc []models.FileInfo
			err := s.filemanager.ReadDir(models.FileArg{}, &dc)
			if err != nil {
				continue
			}
			refreshFunc(dc)
		}
	}()
	// Block until an event is received.
}
