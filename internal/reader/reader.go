package reader

import (
	"log"
	"os"
	"path/filepath"
)

type File struct {
	Name string
	Info os.FileInfo
}

type ListFiles struct {
	Files []File
}

func (l *ListFiles) Walk(path string) (*ListFiles, error) {
	entites, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	file := make([]File, 0, len(entites))
	output := ListFiles{Files: file}

	for _, e := range entites {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			log.Print(err)
			continue
		}
		f := File{
			Name: filepath.Join(path, e.Name()),
			Info: info,
		}
		output.Files = append(output.Files, f)
	}
	return &output, nil
}

func (l *ListFiles) PrintData() {
	for _, file := range l.Files {
		log.Printf("%s, %d, %s, %s\n", file.Name, file.Info.Size(), file.Info.Mode(), file.Info.ModTime())
	}
}

func NewWalker() *ListFiles {
	return &ListFiles{Files: make([]File, 0)}
}
