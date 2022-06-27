package models

type FsMod string

const (
	LOCALFM  FsMod = "local"
	ZIPFM    FsMod = "zip"
	TARFM    FsMod = "tar"
	SYSTRASH FsMod = "trash"
	SFTP     FsMod = "sftp"
)

func IsArchive(mod FsMod) bool {
	return mod == ZIPFM || mod == TARFM
}
