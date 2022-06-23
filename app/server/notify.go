package server

import (
	"atfm/app/models"
	"log"
	"path"

	"github.com/rjeczalik/notify"
)

type NotifyManager struct {
	listeners map[string]chan notify.EventInfo

	onRefresh func(string, []models.FileInfo, bool) error

	filemanager *FileManager
}

func NewNotifyManager(manager *FileManager, onRefresh func(string, []models.FileInfo, bool) error) *NotifyManager {
	return &NotifyManager{
		listeners:   map[string]chan notify.EventInfo{},
		filemanager: manager,
		onRefresh:   onRefresh,
	}
}

func (s *NotifyManager) UnsubscribeRefresh(arg models.FileArg) {
	wp := path.Join(arg.BasePath, arg.Path)
	pc, ok := s.listeners[wp]
	if ok {
		notify.Stop(pc)
		delete(s.listeners, wp)
	}
}

func (s *NotifyManager) SubscribeRefresh(arg models.FileArg) {
	wp := path.Join(arg.BasePath, arg.Path)
	_, ok := s.listeners[wp]
	if ok {
		return
	}
	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch(wp, c, notify.All); err != nil {
		log.Fatal(err)
	}
	s.listeners[wp] = c

	go func() {
		defer notify.Stop(c)
		for {
			ei := <-c
			if ei.Path() == wp {
				err := s.onRefresh(wp, []models.FileInfo{}, true)
				if err != nil {
					return
				}
				continue
			}
			var dc []models.FileInfo
			err := s.filemanager.ReadDir(models.FileArg{
				Mod:      arg.Mod,
				Path:     arg.Path,
				BasePath: arg.BasePath,
			}, &dc)
			if err != nil {
				return
			}
			err = s.onRefresh(wp, dc, ei.Path() == wp)
			if err != nil {
				return
			}
		}
	}()
}
