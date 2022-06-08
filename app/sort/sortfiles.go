package sort

import (
	"atfm/app/models"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

type SortMode int

const (
	SORT_NAME SortMode = iota
	SORT_TYPE
	SORT_SIZE
	SORT_EDITED_DATE
	SORT_RIGHTS
)

type SortPref struct {
	sortMode SortMode
	sortDirs bool
	inverse  bool
}

func NewSortPref() SortPref {
	return SortPref{
		sortMode: SORT_NAME,
		sortDirs: true,
		inverse:  false,
	}
}

func SortContentByName(content []models.FileInfo) []models.FileInfo {
	co := content
	sort.SliceStable(co, func(i, j int) bool {
		return strings.Map(unicode.ToUpper, co[i].Name) <
			strings.Map(unicode.ToUpper, co[j].Name)
	})
	return co
}

func SortContentByType(content []models.FileInfo) []models.FileInfo {
	co := content
	sort.SliceStable(co, func(i, j int) bool {
		extI := filepath.Ext(co[i].Name)
		extJ := filepath.Ext(co[j].Name)
		return strings.Map(unicode.ToUpper, extI) <
			strings.Map(unicode.ToUpper, extJ)
	})
	return co
}

func SortContentBySize(content []models.FileInfo) []models.FileInfo {
	co := content
	sort.SliceStable(co, func(i, j int) bool {
		return co[i].Size > co[j].Size
	})
	return co
}

func InverseContent(content []models.FileInfo) []models.FileInfo {
	co := content
	for i, j := 0, len(co)-1; i < j; i, j = i+1, j-1 {
		co[i], co[j] = co[j], co[i]
	}
	return co
}

func GetContentDirsAndFiles(content []models.FileInfo) ([]models.FileInfo, []models.FileInfo) {
	dirs := make([]models.FileInfo, 0)
	files := make([]models.FileInfo, 0)

	for _, fi := range content {
		if fi.IsDir {
			dirs = append(dirs, fi)
		} else {
			files = append(files, fi)
		}
	}
	return dirs, files
}

func SortDirContent(content []models.FileInfo, prefs SortPref) []models.FileInfo {

	sortC := func(content []models.FileInfo, isDirs bool) []models.FileInfo {
		co := content

		switch prefs.sortMode {
		case SORT_NAME:
			co = SortContentByName(co)

		case SORT_TYPE:
			if isDirs {
				co = SortContentByName(co)
			} else {
				co = SortContentByType(co)
			}

		case SORT_SIZE:
			co = SortContentBySize(co)

		case SORT_EDITED_DATE:

		case SORT_RIGHTS:

		default:
			co = SortContentByName(co)
		}

		if prefs.inverse {
			co = InverseContent(co)
		}
		return co
	}

	if prefs.sortDirs {
		dirs, files := GetContentDirsAndFiles(content)

		dirs = sortC(dirs, true)
		files = sortC(files, false)
		return append(dirs, files...)
	} else {
		return sortC(content, false)
	}

}
