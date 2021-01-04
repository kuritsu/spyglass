package sgc

import (
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/kuritsu/spyglass/api/types"
)

// FileConfig for reading configuration
type FileConfig struct {
	Monitors []*types.Monitor `hcl:"monitor,block"`
	Targets  []*types.Target  `hcl:"target,block"`
}

// File represents a Spyglass Resource Configuration file.
type File struct {
	File     os.FileInfo
	FullPath string
	Config   *FileConfig
}

// ReadConfig using Hashicorp's HCL parser
func (f *File) ReadConfig() error {
	cfg := FileConfig{}
	result := hclsimple.DecodeFile(f.FullPath, nil, &cfg)
	if result != nil {
		return result
	}
	f.Config = &cfg
	return nil
}
