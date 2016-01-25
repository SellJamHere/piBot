package thermo

import (
	"io"
	"os"
)

// Interface file system calls
type fileSystem interface {
	Open(name string) (file, error)
}

// Concrete implementation
type osFS struct{}

func (osFS) Open(name string) (file, error) { return os.Open(name) }

type file interface {
	io.Closer
	io.Reader
}
