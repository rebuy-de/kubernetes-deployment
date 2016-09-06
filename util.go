package main

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

// yes, golang doesn't have any built-in function to copy a file ...
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func FindFiles(dir string, globs ...string) ([]string, error) {
	dir = path.Clean(dir) + "/"
	result := []string{}

	for _, ext := range globs {
		matches, err := filepath.Glob(path.Join(dir, ext))
		if err != nil {
			return nil, err
		}
		result = append(result, matches...)
	}

	return result, nil
}
