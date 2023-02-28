package reader

import (
	"os"

	"github.com/sirupsen/logrus"
)

type File struct {
	Name string      //full path
	Hash uint64      //hash by content
	Info os.FileInfo //file info
}

type ListFiles struct {
	Basepath string
	MapHash  map[uint64]File
	MapName  map[string]File
}

var logger = logrus.New()
var entry *logrus.Entry
