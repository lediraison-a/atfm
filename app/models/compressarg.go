package models

type CompressArg struct {
	Sources     []string
	Destination string
	Mod         FsMod
	BasePath    string
}
