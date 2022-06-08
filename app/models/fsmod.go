package models

type FsMod string

const (
	LOCALFM   FsMod = "local"
	ARCHIVEFM FsMod = "archive"
	SYSTRASH  FsMod = "trash"
	SFTP      FsMod = "sftp"
)
