package models

type ExtractArg struct {
	Source, Destination string
	Mod                 FsMod
	BasePath            string
	OnFileExistAction   FileExistAction
}
