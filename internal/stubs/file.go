package stubs

import (
	"io/fs"
	"os"
)

type stubOS struct {
	OpenFile func(name string, flags int, perm fs.FileMode) (*os.File, error)
}

var OS = stubOS{
	OpenFile: func(name string, flags int, perm fs.FileMode) (*os.File, error) {
		return os.OpenFile(name, flags, perm)
	},
}
