package models

type FileRenameArg struct {
	Mod                     FsMod
	BasePath, Path, NewName string
}
