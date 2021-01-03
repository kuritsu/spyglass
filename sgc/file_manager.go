package sgc

import (
	"io/ioutil"
	"strings"
)

// FileManager implements the Manager using the actual filesystem.
type FileManager struct {
	FileList []File
}

// GetFiles in the directory specified, recursively or not.
func (manager *FileManager) GetFiles(dir string, recursive bool) ([]File, error) {
	list, err := manager.readFiles(dir, recursive)
	manager.FileList = list
	return manager.FileList, err
}

// readFiles in the directory specified, recursively or not.
func (manager *FileManager) readFiles(dir string, recursive bool) ([]File, error) {
	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := []File{}
	for _, f := range list {
		if !f.IsDir() {
			fname := strings.ToLower(f.Name())
			if string(fname[len(fname)-4:]) == ".sgc" {
				result = append(result, File{
					File: f,
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
