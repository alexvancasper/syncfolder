package reader

import (
	"io/fs"
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

func (l *ListFiles) WalkDir(path string) (*ListFiles, error) {
	var output []string

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Print(err)
			return err
		}
		output = append(output, path)
		return nil
	})
	file := make([]File, 0, len(output))
	result := ListFiles{Files: file}

	for i := 0; i < len(output); i++ {
		info, err := os.Stat(output[i])
		if err != nil {
			// log.Print(err)
			continue
		}
		f := File{
			Name: output[i],
			Info: info,
		}

		result.Files = append(result.Files, f)
	}
	return &result, nil
}

func (l *ListFiles) PrintData() {
	for _, file := range l.Files {
		log.Printf("%s, %d, %s, %d\n", file.Name, file.Info.Size(), file.Info.Mode(), file.Info.ModTime().UnixNano())
	}
}

func NewWalker() *ListFiles {
	return &ListFiles{Files: make([]File, 0)}
}
