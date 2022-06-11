package server

import (
	"atfm/app/models"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type XDGTrashManager struct {
	TrashDir    string
	filemanager *FileManager
}

func (t *XDGTrashManager) RestoreFile(file string) (os.FileInfo, string, error) {
	bp := path.Join(t.TrashDir, "info", file+".trashinfo")
	fs := filemanager.GetFs(models.LOCALFM, config.AppConfig.Start.StartBasepath)
	co, _, err := filemanager.OpenFileContent(bp, *fs)
	if err != nil {
		return nil, "", err
	}
	s := strings.Split(string(co), "\n")
	p := ""
	for _, v := range s {
		if strings.HasPrefix(v, "Path=") {
			p = v[strings.Index(v, "=")+1:]
		}
	}
	op := path.Join(t.TrashDir, "files", file)
	errr, c := filemanager.RenameFile(op, p, fs)
	if errr != nil {
		return nil, "", errr
	}
	err = filemanager.DeleteFile(path.Join(bp, file), *fs)
	if err != nil {
		return nil, "", err
	}
	return c, p, errr
}

func (t *XDGTrashManager) TrashFile(file string) error {
	fs := filemanager.GetFs(filemanager.LOCALFM, config.AppConfig.Start.StartBasepath)
	nf := path.Join(t.TrashDir, "info", filepath.Base(file)+".trashinfo")
	np := path.Join(t.TrashDir, "files", filepath.Base(file))
	exist, err := filemanager.Exist(np, *fs)
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
		exist, err = filemanager.Exist(np, *fs)
		if err != nil {
			return err
		}
	}
	_, f, err := filemanager.CreateFile(nf, *fs)
	if err != nil {
		return err
	}
	dte := time.Now().Format("2022-04-16T17:22:22")
	trashInfo := "[Trash Info]\n"
	trashInfo += "Path=" + file + "\n"
	trashInfo += "DeletionDate=" + dte + "\n"
	_, errrr := f.WriteString(trashInfo)
	if errrr != nil {
		return errrr
	}
	errr, _ := filemanager.RenameFile(file, np, fs)
	if errr != nil {
		return errr
	}
	return nil
}

func (t *XDGTrashManager) GetTrashInfo(file string) (string, string, error) {
	bp := path.Join(t.TrashDir, "info")
	fs := filemanager.GetFs(filemanager.LOCALFM, config.AppConfig.Start.StartBasepath)
	co, _, err := filemanager.OpenFileContent(path.Join(bp, file), *fs)
	if err != nil {
		return "", "", err
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
	return p, d, nil
}

func (t *XDGTrashManager) EmptyTrash() {
	fs := filemanager.GetFs(models.LOCALFM, config.AppConfig.Start.StartBasepath)
	filesPath := path.Join(t.TrashDir, "files")
	infoPath := path.Join(t.TrashDir, "info")
	c := filemanager.GetDirsInfoAt(filesPath, fs)
	if c != nil {
		for _, fi := range c {
			filemanager.DeleteFile(path.Join(filesPath, fi.Name()), *fs)
		}
	}
	c = filemanager.GetDirsInfoAt(infoPath, fs)
	if c == nil {
		return
	}
	for _, fi := range c {
		if filepath.Ext(fi.Name()) == ".trashinfo" {
			filemanager.DeleteFile(path.Join(infoPath, fi.Name()), *fs)
		}
	}
}
