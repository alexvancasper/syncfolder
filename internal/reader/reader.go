package reader

import (
	"strings"
)

func NewWalker() *ListFiles {
	return &ListFiles{Basepath: "", MapHash: make(map[uint64]File), MapName: make(map[string]File)}
}

func MakeNameKey(path, filepath string) string {
	res, _ := strings.CutPrefix(filepath, path)
	return res
}
