package server

import (
	"atfm/app/models"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

var testBasePath = "/home/alban/go/src/atfm/tests"

func TestNewFileManager(t *testing.T) {
	fm := NewFileManager()
	if fm == nil {
		t.Fatal("Filemanager is null")
	}
	if fm.filesystems == nil {
		t.Fatal("Filemanager.filesystems is null")
	}
}

func TestGetFs(t *testing.T) {
	fm := NewFileManager()
	fsLocal, errLocal := fm.getFs("/", models.LOCALFM)
	if errLocal != nil {
		t.Fatal(errLocal)
	}

	if afs, ok := fsLocal.(afero.OsFs); ok {
		t.Fatalf("Filemanager getFs dont give good LOCAL afero fs, want afero.OsFs, get %s", reflect.TypeOf(afs))
	}

	fsBp, errBp := fm.getFs(testBasePath, models.LOCALFM)
	if errBp != nil {
		t.Fatal(errBp)
	}

	if fsBp == nil {
		t.Fatal("Filemanager getFs with basepath != '/' give null afero.Fs")
	}
}

func TestReadDir(t *testing.T) {
	fm := NewFileManager()
	arg := models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: testBasePath,
		Path:     "/",
	}
	var dc []models.FileInfo
	err := fm.ReadDir(arg, &dc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStatFile(t *testing.T) {
	fm := NewFileManager()
	arg := models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: testBasePath,
		Path:     "/dir",
	}
	var di models.FileInfo
	err := fm.StatFile(arg, &di)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRenameFile(t *testing.T) {
	fm := NewFileManager()
	arg := models.FileRenameArg{
		Mod:      models.LOCALFM,
		BasePath: testBasePath,
		Path:     "/dir/file.txt",
		NewName:  "/dir/file2.txt",
	}
	var di models.FileInfo
	err := fm.RenameFile(arg, &di)
	if err != nil {
		t.Fatal(err)
	}
	fm.RenameFile(models.FileRenameArg{
		Mod:      models.LOCALFM,
		BasePath: testBasePath,
		Path:     "/dir/file2.txt",
		NewName:  "/dir/file.txt",
	}, &di)
}

func TestExist(t *testing.T) {
	fm := NewFileManager()
	arg := models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: testBasePath,
		Path:     "/dir",
	}
	var ex bool
	err := fm.Exist(arg, &ex)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateFile(t *testing.T) {
	fm := NewFileManager()
	arg := models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: testBasePath,
		Path:     "/dir/file2.txt",
	}
	var fi models.FileInfo
	err := fm.CreateFile(arg, &fi)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadFile(t *testing.T) {
	fm := NewFileManager()
	arg := models.FileArg{
		Mod:      models.LOCALFM,
		BasePath: testBasePath,
		Path:     "/dir/file.txt",
	}
	var fc []byte
	err := fm.ReadFile(arg, &fc)
	if err != nil {
		t.Fatal(err)
	}
}
