package config

import (
	"fmt"
	"io/fs"
	"os"
)

type MachineFS struct{}

func (fs MachineFS) Open(name string) (fs.File, error) {
	return nil, fmt.Errorf("not implemented")
}

func (fs MachineFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
