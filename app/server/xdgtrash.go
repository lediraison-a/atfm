package server

import (
	"atfm/app/models"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type XDGTrashManager struct {
	TrashDir string
	fm       *FileManager
}

func NewXDGTrashManager(trashDir string, fileManager *FileManager) (*XDGTrashManager, error) {
	filesPath := path.Join(trashDir, "files")
	infoPath := path.Join(trashDir, "info")
	var fpEx, ipEx bool
	err := fileManager.Exist(models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: "/",
		Path:     filesPath,
	}, &fpEx)
	if err != nil {
		return nil, err
	}
	err = fileManager.Exist(models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: "/",
		Path:     infoPath,
	}, &ipEx)
	if err != nil {
		return nil, err
	}
	if !fpEx || !ipEx {
		return nil, errors.New("xdg trash directories doesn't exists")
	}

	return &XDGTrashManager{
		TrashDir: trashDir,
		fm:       fileManager,
	}, nil
}

func (t *XDGTrashManager) RestoreFile(file string, fileInfo *models.FileInfo) error {
	bp := path.Join(t.TrashDir, "info", file+".trashinfo")
	var co []byte
	err := t.fm.ReadFile(models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: "/",
		Path:     bp,
	}, &co)
	if err != nil {
		return err
	}
	s := strings.Split(string(co), "\n")
	p := ""
	for _, v := range s {
		if strings.HasPrefix(v, "Path=") {
			p = v[strings.Index(v, "=")+1:]
		}
	}
	op := path.Join(t.TrashDir, "files", file)
	var c models.FileInfo
	err = t.fm.RenameFile(models.FileRenameArg{
		Mod:      models.LOCALFM,
		BasePath: "/",
		Path:     op,
		NewName:  p,
	}, &c)
	if err != nil {
		return err
	}
	err = t.fm.DeleteFile(models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: "/",
		Path:     bp,
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *XDGTrashManager) TrashFile(file string, fileInfo *models.FileInfo) error {
	nf := path.Join(t.TrashDir, "info", filepath.Base(file)+".trashinfo")
	np := path.Join(t.TrashDir, "files", filepath.Base(file))
	var exist bool
	err := t.fm.Exist(models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: "/",
		Path:     np,
	}, &exist)
	if err != nil {
		return err
	}
	i := 0
	for exist {
		ext := filepath.Ext(file)
		filename := filepath.Base(file)
		name := filename[0 : len(filename)-len(ext)]
		np = path.Join(t.TrashDir, "files", fmt.Sprintf("%s%d%s", name, i, ext))
		nf = path.Join(t.TrashDir, "info", fmt.Sprintf("%s%d%s", name, i, ext)+".trashinfo")
		err := t.fm.Exist(models.FileArg{
			Mod:      models.LOCALFM,
			BasePath: "/",
			Path:     np,
		}, &exist)
		if err != nil {
			return err
		}
	}
	var f models.FileInfo
	dte := time.Now().Format("2022-04-16T17:22:22")
	trashInfo := "[Trash Info]\n"
	trashInfo += "Path=" + file + "\n"
	trashInfo += "DeletionDate=" + dte + "\n"
	err = t.fm.CreateAndWriteFile(models.FileWriteArg{
		Mod:      models.LOCALFM,
		BasePath: "/",
		Path:     nf,
		Content:  trashInfo,
	}, &f)
	if err != nil {
		return err
	}
	var ft models.FileInfo
	err = t.fm.RenameFile(models.FileRenameArg{
		Mod:      models.LOCALFM,
		BasePath: "/",
		Path:     file,
		NewName:  np,
	}, &ft)
	if err != nil {
		return err
	}
	*fileInfo = ft
	return nil
}

func (t *XDGTrashManager) GetTrashInfo(file string, trashInfo *models.TrashInfo) error {
	bp := path.Join(t.TrashDir, "info")
	fibp := path.Join(bp, file+".trashinfo")
	var co []byte
	err := t.fm.ReadFile(models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: "/",
		Path:     fibp,
	}, &co)
	if err != nil {
		return err
	}
	s := strings.Split(string(co), "\n")
	p, d := "", ""
	for _, v := range s {
		if strings.HasPrefix(v, "Path=") {
			p = v[strings.Index(v, "="):]
		} else if strings.HasPrefix(v, "DeletionDate=") {
			d = v[strings.Index(v, "="):]
		}
	}
	if p == "" || d == "" {
		return errors.New("")
	}
	*trashInfo = models.TrashInfo{
		Path:        fibp,
		RestorePath: p,
		TrashDate:   d,
	}
	return nil
}

func (t *XDGTrashManager) EmptyTrash() error {
	filesPath := path.Join(t.TrashDir, "files")
	infoPath := path.Join(t.TrashDir, "info")
	var c []models.FileInfo
	err := t.fm.ReadDir(models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: "/",
		Path:     path.Join(t.TrashDir, filesPath),
	}, &c)
	if err != nil {
		return err
	}
	for _, fi := range c {
		var ti models.TrashInfo
		err := t.GetTrashInfo(fi.Name, &ti)
		if err != nil {
			continue
		}
		t.fm.DeleteFile(models.FileArg{
			Mod:      models.LOCALFM,
			BasePath: "/",
			Path:     ti.Path,
		})

		if err != nil {
			continue
		}
		err = t.fm.DeleteFile(models.FileArg{
			Mod:      models.LOCALFM,
			BasePath: "/",
			Path:     path.Join(infoPath, fi.Name+".trashinfo"),
		})
		if err != nil {
			continue
		}
	}

	return nil
}
