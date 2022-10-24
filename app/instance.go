package app

import (
	"atfm/app/config"
	"atfm/app/models"
	"atfm/app/server"
	"atfm/app/sort"
	"atfm/generics"
	"errors"
	"net/rpc"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var FileManagerService *server.FileManager

var NotifyManagerService *server.NotifyManager

type Instance struct {
	rpcClient *rpc.Client

	Id int

	History *models.NavHistory

	BasePath string
	Mod      models.FsMod

	DirPath string
	DirInfo models.FileInfo

	Content      []models.FileInfo
	ShownContent []models.FileInfo

	SelectedIndexes []int

	ShowHidden, ShowOpenParent bool

	CurrentItem int

	QuickSearch *Search

	config *config.Config
}

func NewInstance(mod models.FsMod, path, basePath string, id int, showOpenParent bool) *Instance {
	return &Instance{
		Id:              id,
		History:         models.NewNavHistory(),
		Mod:             mod,
		DirInfo:         models.FileInfo{},
		Content:         []models.FileInfo{},
		ShownContent:    []models.FileInfo{},
		SelectedIndexes: []int{},
		DirPath:         path,
		ShowHidden:      false,
		CurrentItem:     0,
		BasePath:        basePath,
		QuickSearch:     NewSearch(),
		ShowOpenParent:  showOpenParent,
	}
}

func (s *Instance) FileCount() int {
	return len(s.ShownContent)
}

func (s *Instance) OpenAtIndex(index int) error {
	if index < 0 || index > len(s.ShownContent)-1 {
		return errors.New("index out of range")
	}
	p := path.Join(s.DirPath, s.ShownContent[index].Name)
	return s.OpenDirSaveHistory(p, s.BasePath, s.Mod)
}

func (s *Instance) OpenDirSaveHistory(path, basepath string, mod models.FsMod) error {
	h := s.GetHistoryRecCurrent()
	err := s.OpenDir(path, basepath, mod)
	if err == nil {
		s.History.GetHistoryStack(models.HISTORY_BACK).Push(&h)
	}
	return err
}

func (s *Instance) OpenDir(path, basepath string, mod models.FsMod) error {
	if path == "" || basepath == "" {
		return nil
	}
	if path[0] == '$' {
		pathEnv := os.Getenv(path[1:])
		if pathEnv == "" {
			return nil
		}
		path = pathEnv
	}

	if filepath.Ext(path) == ".zip" {
		basepath = path
		path = "/"
		mod = models.ZIPFM
	} else if strings.HasSuffix(path, ".tar.gz") ||
		strings.HasSuffix(path, ".tar.xz") {
		basepath = path
		path = "/"
		mod = models.TARFM
	}

	if path == "/" && s.DirPath == "/" && (s.Mod == models.TARFM || s.Mod == models.ZIPFM) {
		path = filepath.Dir(s.BasePath)
		basepath = "/"
		mod = models.LOCALFM
	}

	arg := models.FileArg{
		Mod:      mod,
		BasePath: basepath,
		Path:     path,
	}

	var dc []models.FileInfo
	// err := s.rpcClient.Call("FileManager.ReadDir", arg, &dc)
	err := FileManagerService.ReadDir(arg, &dc)
	if err != nil {
		return err
	}
	var di models.FileInfo
	// err = s.rpcClient.Call("FileManager.StatFile", arg, &di)
	err = FileManagerService.StatFile(arg, &di)
	if err != nil {
		return err
	}

	if mod == models.LOCALFM {
		NotifyManagerService.UnsubscribeRefresh(models.FileArg{
			Mod:      s.Mod,
			BasePath: s.BasePath,
			Path:     s.DirPath,
		})
		NotifyManagerService.SubscribeRefresh(arg)
	}

	s.BasePath = basepath
	s.Mod = mod
	s.DirPath = path
	s.Content = dc
	s.ShownContent = s.GetShownContent(dc)
	s.DirInfo = di
	s.CurrentItem = 0
	s.QuickSearch.ResetSearch()
	return nil
}

func (s *Instance) RenameFile(filepath, newname string) error {
	arg := models.FileRenameArg{
		Mod:      s.Mod,
		BasePath: s.BasePath,
		Path:     filepath,
		NewName:  newname,
	}
	var fi models.FileInfo
	return FileManagerService.RenameFile(arg, &fi)
}

func (s *Instance) NewFile(filepath string) error {
	arg := models.FileArg{
		Mod:      s.Mod,
		BasePath: s.BasePath,
		Path:     filepath,
	}
	var fi models.FileInfo
	return FileManagerService.CreateFile(arg, &fi)
}

func (s Instance) NewDir(filepath string) error {
	arg := models.FileArg{
		Mod:      s.Mod,
		BasePath: s.BasePath,
		Path:     filepath,
	}
	var fi models.FileInfo
	return FileManagerService.CreateDir(arg, &fi)
}

func (s *Instance) ReadDir(path, basepath string, mod models.FsMod) ([]models.FileInfo, error) {
	arg := models.FileArg{
		Mod:      mod,
		BasePath: basepath,
		Path:     path,
	}
	var dc []models.FileInfo
	// err := s.rpcClient.Call("FileManager.ReadDir", arg, &dc)
	err := FileManagerService.ReadDir(arg, &dc)
	if err != nil {
		return nil, err
	}
	return s.GetShownContent(dc), nil
}

func (s *Instance) OpenHistoryDir(mod models.NavHistoryMod) (bool, error) {
	hStack := s.History.GetHistoryStack(mod)
	if hStack.Count > 0 {
		h := hStack.Pop()
		err := s.OpenDir(h.Path, h.BasePath, h.Mod)
		if err == nil {
			s.CurrentItem = h.Index
			s.History.GetHistoryStack(!mod).Push(h)
			return true, nil
		} else {
			return false, err
		}
	} else {
		return false, nil
	}
}

func (s *Instance) IsSelected(index int) bool {
	return generics.Contains(s.SelectedIndexes, index)
}

func (s *Instance) CanShowOpenParent() bool {
	return s.ShowOpenParent &&
		(s.DirPath != "/" || (s.Mod == models.TARFM || s.Mod == models.ZIPFM))
}

func (s *Instance) GetHistoryRecCurrent() models.NavHistoryRec {
	return models.NavHistoryRec{
		Path:     s.DirPath,
		Index:    s.CurrentItem,
		Mod:      s.Mod,
		BasePath: s.BasePath,
	}
}

func (s *Instance) GetParentInfo() (string, string, models.FsMod) {
	pth := s.DirPath
	bp := s.BasePath
	mod := s.Mod
	if pth == "/" && models.IsArchive(s.Mod) {

	} else {
		pth = path.Dir(s.DirPath)
	}
	return pth, bp, mod
}

func (s *Instance) SelectItem(index int, toggle bool) (int, bool) {
	if index < 0 {
		return 0, false
	}
	if index == 0 && s.CanShowOpenParent() {
		return 0, false
	}
	for i, v := range s.SelectedIndexes {
		if v == index {
			if toggle {
				s.SelectedIndexes = append(s.SelectedIndexes[:i], s.SelectedIndexes[i+1:]...)
				return v, false
			}
			return v, true
		}
	}
	s.SelectedIndexes = append(s.SelectedIndexes, index)
	return index, true
}

func (s *Instance) UnselectAll() {
	s.SelectedIndexes = []int{}
}

func (s *Instance) GetShownContent(content []models.FileInfo) []models.FileInfo {
	sdc := sort.SortDirContent(content, sort.NewSortPref())
	if !s.ShowHidden {
		sdc = generics.Filter(sdc, func(v models.FileInfo, index int) bool {
			return !strings.HasPrefix(v.Name, ".")
		})
	}
	if s.CanShowOpenParent() {
		parent := models.FileInfo{
			Name:    "..",
			IsDir:   true,
			Mode:    0,
			Size:    0,
			ModTime: time.Time{},
			Symlink: "",
		}
		sdc = append([]models.FileInfo{parent}, sdc...)
	}
	return sdc
}

func (s *Instance) IsEmpty() bool {
	return len(s.ShownContent) == 0
}
