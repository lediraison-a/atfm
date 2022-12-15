package models

type FileExistAction string

const (
	FILE_REPLACE     FileExistAction = "replace"
	FILE_SKIP        FileExistAction = "skip"
	OPERATION_CANCEL FileExistAction = "cancel"
)
