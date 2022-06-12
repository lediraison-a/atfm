package server

import (
	"github.com/rjeczalik/notify"
	"log"
)

type NotifyManager struct {
}

func (s *NotifyManager) SuscribeNotify(instanceId int, path string) {
	c := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening on events within current working directory.
	// Dispatch each create and remove events separately to c.
	if err := notify.Watch(path, c, notify.Create, notify.Remove); err != nil {
		log.Fatal(err)
	}

	go func() {
		defer notify.Stop(c)
		for {
			ei := <-c
			log.Println("Got event:", ei)
		}
	}()
	// Block until an event is received.
}
