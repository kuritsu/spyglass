package sgc

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// FileManager implements the Manager using the actual filesystem.
type FileManager struct {
	FileList []*File
}

// GetFiles in the directory specified, recursively or not.
func (manager *FileManager) GetFiles(dir string, recursive bool) ([]*File, error) {
	list, err := manager.readFiles(dir, recursive)
	manager.FileList = list
	return manager.FileList, err
}

// readFiles in the directory specified, recursively or not.
func (manager *FileManager) readFiles(dir string, recursive bool) ([]*File, error) {
	fpath, perr := filepath.Abs(dir)
	if perr != nil {
		return nil, perr
	}
	list, err := ioutil.ReadDir(fpath)
	if err != nil {
		return nil, err
	}
	result := []*File{}
	for _, f := range list {
		if !f.IsDir() {
			fname := strings.ToLower(f.Name())
			if string(fname[len(fname)-4:]) == ".hcl" {
				result = append(result, &File{
					File:     f,
					FullPath: fpath + string(filepath.Separator) + f.Name(),
				})
			}
		} else if recursive {
			subdirList, serr := manager.readFiles(dir+"/"+f.Name(), recursive)
			if serr != nil {
				return nil, serr
			}
			for _, s := range subdirList {
				result = append(result, s)
			}
		}
	}
	return result, nil
}

// ParseConfig from file list
func (manager *FileManager) ParseConfig() []FileParseError {
	if manager.FileList == nil {
		return []FileParseError{{
			File:       nil,
			InnerError: errors.New("No files read"),
		}}
	}
	result := []FileParseError{}
	for _, f := range manager.FileList {
		err := f.ReadConfig()
		if err != nil {
			result = append(result, FileParseError{File: f, InnerError: err})
		}
	}
	if len(result) > 0 {
		return result
	}
	return nil
}
