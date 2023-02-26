package types

import (
	"os"
)

type File struct {
	Name string
	Info os.FileInfo
}

type ListFiles []File
