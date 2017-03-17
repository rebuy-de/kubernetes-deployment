package cmd

import (
	"io"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
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

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func Must(err error) {
	if err == nil {
		return
	}

	log.Error(err)

	if err, ok := err.(stackTracer); ok {
		log.Debugf("%+v", err.StackTrace())
	}

	os.Exit(1)
}
