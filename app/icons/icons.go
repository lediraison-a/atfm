package icons

import "path/filepath"

type Icon struct {
	icon  string
	color string
	name  string
}

func GetFileIcon(name string, colorDefault string) (string, string) {
	fnicon, ok := FileNameIcons[name]
	if ok {
		return fnicon, colorDefault
	}
	fticon, ok := FileTypeIcons[filepath.Ext(name)]
	if ok {
		return fticon.icon, fticon.color
	}
	return FileIcon, colorDefault
}

func GetDirIcon(name string) string {
	ic, ok := FileNameIcons[name]
	if ok {
		return ic
	}
	return DirIcon
}

var FileIcon = ""

var DirIcon = ""
