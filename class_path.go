package gj

import (
	"archive/zip"
	"fmt"
	"github.com/murakmii/gj/class_file"
	"os"
	"path/filepath"
	"strings"
)

type (
	ClassPath interface {
		SearchClass(name string) (*class_file.Class, error)
		Close()
	}

	jar struct {
		r *zip.ReadCloser
	}

	dir struct {
		path string
	}
)

func InitClassPaths(paths []string) (classPaths []ClassPath, err error) {
	classPaths = make([]ClassPath, 0)
	defer func() {
		if err != nil {
			for _, cp := range classPaths {
				cp.Close()
			}
			classPaths = nil
		}
	}()

	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			return
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				return
			}

			var cp ClassPath
			if info.IsDir() {
				cp = &dir{path: match}

			} else if strings.HasSuffix(match, ".jar") {
				j := &jar{}
				j.r, err = zip.OpenReader(match)
				if err != nil {
					return
				}
				cp = j

			} else {
				return nil, fmt.Errorf("unsupported class path entry: %s", match)
			}

			classPaths = append(classPaths, cp)
		}
	}

	return
}

func (j *jar) SearchClass(name string) (*class_file.Class, error) {
	cfReader, err := j.r.Open(name)
	if err != nil {
		return nil, err
	}
	defer cfReader.Close()

	return class_file.ReadClassFile(cfReader)
}

func (j *jar) Close() {
	j.r.Close()
}

func (d *dir) SearchClass(name string) (*class_file.Class, error) {
	return class_file.OpenClassFile(filepath.Join(d.path, name))
}

func (d *dir) Close() {
	// do nothing
}
