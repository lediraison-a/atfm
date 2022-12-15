package server

import (
	"archive/zip"
	"atfm/app/models"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// https://gist.github.com/yhirose/addb8d248825d373095c

func ExtractZip(source, destination string, destFs afero.Fs, onFileExist models.FileExistAction) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := destFs.MkdirAll(destination, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(destination, file.Name)
		var exist bool
		exist, err = afero.Exists(destFs, path)
		if err != nil {
			return err
		}
		if exist {
			switch onFileExist {
			case models.OPERATION_CANCEL:
				return nil
			case models.FILE_SKIP:
				continue
			case models.FILE_REPLACE:
			}
		}

		if file.FileInfo().IsDir() {
			err := destFs.MkdirAll(path, file.Mode())
			if err != nil {
				return err
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := destFs.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func CompressZip(sources []string, destination string, needBaseDir bool, destFs afero.Fs) error {
	zipfile, err := destFs.Create(destination)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	for _, source := range sources {
		info, err := destFs.Stat(source)
		if err != nil {
			return err
		}

		var baseDir string
		if info.IsDir() {
			baseDir = filepath.Base(source)
		}

		err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			if baseDir != "" {
				if needBaseDir {
					header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
				} else {
					path := strings.TrimPrefix(path, source)
					if len(path) > 0 && (path[0] == '/' || path[0] == '\\') {
						path = path[1:]
					}
					if len(path) == 0 {
						return nil
					}
					header.Name = path
				}
			}

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := destFs.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		})
		if err != nil {
			return err
		}
	}

	return err
}
