package models

type FileWriteArg struct {
	Mod            FsMod
	BasePath, Path string
}
