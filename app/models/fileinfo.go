package models

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	Name    string
	IsDir   bool
	Mode    fs.FileMode
	Size    int64
	ModTime time.Time
}
