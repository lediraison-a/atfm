package server

import (
	"atfm/app/models"
	"testing"
	"time"
)

func TestNewNotifyManager(t *testing.T) {
	// nm := NewNotifyManager()
}

func TestSubscribeRefresh(t *testing.T) {
	p := "/"
	done := false
	sf := false
	nm := NewNotifyManager(NewFileManager(), func(_ string, fi []models.FileInfo, selfDelete bool) error {
		done = true
		if selfDelete {
			sf = true
		}
		t.Logf("files : %d", len(fi))
		return nil
	})
	var fi models.FileInfo

	nm.SubscribeRefresh(models.FileArg{
		Path:     p,
		BasePath: testBasePath,
		Mod:      models.LOCALFM,
	})

	time.Sleep(1 * time.Second)

	fa := models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: testBasePath,
		Path:     "/test.txt",
	}
	err := nm.filemanager.CreateFile(fa, &fi)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	if !done {
		t.Fatal("didn't refreshed")
	}
	nm.UnsubscribeRefresh(fa)
	err = nm.filemanager.DeleteFile(fa)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	p = "/dirtest"
	fa = models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: testBasePath,
		Path:     p,
	}
	err = nm.filemanager.CreateDir(fa, &fi)
	if err != nil {
		t.Fatal(err)
	}
	nm.SubscribeRefresh(fa)
	time.Sleep(1 * time.Second)

	err = nm.filemanager.DeleteFile(fa)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)
	if !sf {
		t.Fatal("didn't self deleted detected")
	}

	nm.UnsubscribeRefresh(fa)
}
