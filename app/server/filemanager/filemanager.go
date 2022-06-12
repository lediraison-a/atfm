package filemanager

import "context"

type FileManager struct {
}

func (s *FileManager) ReadDir(ctx context.Context, arg *FileArg) (*FileInfos, error) {
	return &FileInfos{
		Files: []*FileInfo{},
	}, nil
}
