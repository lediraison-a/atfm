package server

import (
	"atfm/app/models"
	"errors"
	"os"
	"path"

	"github.com/spf13/afero"
)

type FsId struct {
	mod      models.FsMod
	basepath string
}

func NewFsId(basepath string, mod models.FsMod) FsId {
	return FsId{
		mod:      mod,
		basepath: basepath,
	}
}

type FileManager struct {
	filesystems map[FsId]afero.Fs
}

func NewFileManager() *FileManager {
	return &FileManager{
		filesystems: map[FsId]afero.Fs{},
	}
}

func (f *FileManager) getFs(basepath string, mod models.FsMod) (afero.Fs, error) {
	id := NewFsId(basepath, mod)
	dfs, ok := f.filesystems[id]
	if ok {
		return dfs, nil
	}
	switch mod {
	case models.LOCALFM:
		dfs = afero.NewOsFs()
		id := NewFsId(basepath, models.LOCALFM)
		if basepath != "/" {
			bpfs := afero.NewBasePathFs(dfs, basepath)
			if bpfs == nil {
				return nil, errors.New("cannot create afero BasePathFs : wrong basepath")
			}
			f.filesystems[id] = bpfs
			return bpfs, nil
		}
		f.filesystems[id] = dfs
		return dfs, nil
	}
	return nil, errors.New("no file system to open path")
}

func (f *FileManager) ReadDir(arg models.FileArg, dirContent *[]models.FileInfo) error {
	dfs, err := f.getFs(arg.BasePath, arg.Mod)
	if err != nil {
		return err
	}
	fi, err := dfs.Open(arg.Path)
	if err != nil {
		return err
	}
	dc, err := fi.Readdir(0)
	if err != nil {
		return err
	}
	dci := []models.FileInfo{}
	for _, entry := range dc {
		dci = append(dci, models.FileInfo{
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			Mode:    entry.Mode(),
			Size:    entry.Size(),
			ModTime: entry.ModTime(),
			Symlink: f.readSimlink(path.Join(arg.Path, entry.Name()), dfs, arg.Mod, entry),
		})
	}
	*dirContent = dci
	return nil
}

func (f *FileManager) StatFile(arg models.FileArg, fileInfo *models.FileInfo) error {
	dfs, err := f.getFs(arg.BasePath, arg.Mod)
	if err != nil {
		return err
	}
	di, err := dfs.Stat(arg.Path)
	if err != nil {
		return err
	}
	*fileInfo = models.FileInfo{
		Name:    di.Name(),
		IsDir:   di.IsDir(),
		Mode:    di.Mode(),
		Size:    di.Size(),
		ModTime: di.ModTime(),
		Symlink: f.readSimlink(arg.Path, dfs, arg.Mod, di),
	}
	return nil
}

func (f *FileManager) RenameFile(arg models.FileRenameArg, fileInfo *models.FileInfo) error {
	fs, err := f.getFs(arg.BasePath, arg.Mod)
	if err != nil {
		return err
	}
	if arg.Path != arg.NewName {
		err = fs.Rename(arg.Path, arg.NewName)
		if err != nil {
			return err
		}
	}
	var fi models.FileInfo
	err = f.StatFile(models.FileArg{
		Mod:      arg.Mod,
		BasePath: arg.BasePath,
		Path:     arg.NewName,
	}, &fi)
	if err != nil {
		return err
	}
	*fileInfo = fi
	return nil
}

func (f *FileManager) Exist(arg models.FileArg, exist *bool) error {
	fs, err := f.getFs(arg.BasePath, arg.Mod)
	if err != nil {
		return err
	}
	b, err := afero.Exists(fs, arg.Path)
	if err != nil {
		return err
	}
	*exist = b
	return nil
}

func (f *FileManager) CreateFile(arg models.FileArg, fileInfo *models.FileInfo) error {
	fs, err := f.getFs(arg.BasePath, arg.Mod)
	if err != nil {
		return err
	}
	fi, err := fs.Create(arg.Path)
	if err != nil {
		return err
	}
	fii, err := fi.Stat()
	if err != nil {
		return err
	}
	*fileInfo = models.FileInfo{
		Name:    fii.Name(),
		IsDir:   fii.IsDir(),
		Mode:    fii.Mode(),
		Size:    fii.Size(),
		ModTime: fii.ModTime(),
		Symlink: f.readSimlink(arg.Path, fs, arg.Mod, fii),
	}
	return nil
}

func (f *FileManager) CreateAndWriteFile(arg models.FileWriteArg, fileInfo *models.FileInfo) error {
	fs, err := f.getFs(arg.BasePath, arg.Mod)
	if err != nil {
		return err
	}
	fi, err := fs.Create(arg.Path)
	if err != nil {
		return err
	}
	_, err = fi.Write([]byte(arg.Content))
	if err != nil {
		return err
	}
	fii, err := fi.Stat()
	if err != nil {
		return err
	}
	*fileInfo = models.FileInfo{
		Name:    fii.Name(),
		IsDir:   fii.IsDir(),
		Mode:    fii.Mode(),
		Size:    fii.Size(),
		ModTime: fii.ModTime(),
	}
	return nil
}

func (f *FileManager) ReadFile(arg models.FileArg, content *[]byte) error {
	fs, err := f.getFs(arg.BasePath, arg.Mod)
	if err != nil {
		return err
	}
	fc, err := afero.ReadFile(fs, arg.Path)
	if err != nil {
		return err
	}
	*content = fc
	return nil
}

func (f *FileManager) DeleteFile(arg models.FileArg) error {
	fs, err := f.getFs(arg.BasePath, arg.Mod)
	if err != nil {
		return err
	}
	err = fs.Remove(arg.Path)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileManager) WriteReader() error {

	return nil
}

func (f *FileManager) readSimlink(name string, fs afero.Fs, mod models.FsMod, info os.FileInfo) string {
	t := ""
	var err error
	if mod != models.LOCALFM {
		return t
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return t
	}
	if afs, ok := fs.(*afero.OsFs); ok {
		t, err = afs.ReadlinkIfPossible(name)
		if err != nil {
			return t
		}
	} else if afs, ok := fs.(*afero.BasePathFs); ok {
		t, err = afs.ReadlinkIfPossible(name)
		if err != nil {
			return t
		}
	} else if afs, ok := fs.(*afero.ReadOnlyFs); ok {
		t, err = afs.ReadlinkIfPossible(name)
		if err != nil {
			return t
		}
	}

	return t
}
