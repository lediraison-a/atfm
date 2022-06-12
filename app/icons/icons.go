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
	ext := filepath.Ext(name)
	if len(ext) > 1 {
		fticon, ok := FileTypeIcons[ext[1:]]
		if ok {
			return fticon.icon, fticon.color
		}
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

var DirEmptyIcon = ""
