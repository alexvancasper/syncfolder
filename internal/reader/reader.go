package reader

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/zeebo/xxh3"
)

type File struct {
	Name    string
	Hash    uint64
	CumHash uint64
	Info    os.FileInfo
}

type ListFiles struct {
	Files []File
	Map   map[uint64]File
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
	result := ListFiles{Files: file, Map: make(map[uint64]File)}

	for i := 0; i < len(output); i++ {
		info, err := os.Stat(output[i])
		if err != nil {
			// log.Print(err)
			continue
		}
		hashVal := hashCalc(output[i])
		f := File{
			Name:    output[i],
			Hash:    hashVal,
			CumHash: hashCalcCum(fmt.Sprintf("%s%d", filepath.Base(output[i]), hashVal)),
			Info:    info,
		}
		result.Files = append(result.Files, f)
		result.Map[hashVal] = f
	}
	return &result, nil
}

func (l *ListFiles) PrintData() {
	for _, file := range l.Files {
		log.Printf("%s, %d, %d, %d\n", file.Name, file.Hash, file.Info.Size(), file.Info.ModTime().UnixNano())
	}
}

func NewWalker() *ListFiles {
	return &ListFiles{Files: make([]File, 0)}
}

func hashCalc(path string) uint64 {
	fileSize, err := os.Stat(path)
	if err != nil {
		log.Print(err)
	}

	buf := make([]byte, 0, fileSize.Size())

	readFile, err := os.Open(path)
	if err != nil {
		log.Print(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	for fileScanner.Scan() {
		buf = append(buf, fileScanner.Bytes()...)
	}
	return xxh3.Hash(buf)
}

func hashCalcCum(data string) uint64 {
	return xxh3.Hash([]byte(data))
}
