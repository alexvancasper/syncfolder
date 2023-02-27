package sync

import (
	"bufio"
	"final_task/internal/reader"
	"log"
	"os"
	"path/filepath"

	"github.com/zeebo/xxh3"
)

func Intersection(dir1, dir2 reader.ListFiles) reader.ListFiles {
	if len(dir1.Files) == 0 || len(dir2.Files) == 0 {
		return reader.ListFiles{}
	}
	var hashA uint64

	set := make([]reader.File, 0, len(dir1.Files))
	for _, a := range dir1.Files {
		if a.Info.IsDir() {
			continue
		}
		hashA = hashCalc(a.Name)
		for _, b := range dir2.Files {
			if b.Info.IsDir() {
				continue
			}
			hashB := hashCalc(b.Name)
			if filepath.Base(a.Name) == filepath.Base(b.Name) && hashA == hashB {
				set = append(set, a)
				// log.Printf("File: %s hash: %d AND File: %s hash: %d - matched", a.Name, hashA, b.Name, hashB)
			} else {
				// log.Printf("File: %s hash: %d AND File: %s hash: %d - not matched", a.Name, hashA, b.Name, hashB)
			}
		}
	}

	return reader.ListFiles{Files: set}
}

func Distinction(dir1, dir2 reader.ListFiles) reader.ListFiles {
	if len(dir1.Files) == 0 || len(dir2.Files) == 0 {
		return reader.ListFiles{}
	}
	intersec := Intersection(dir1, dir2)
	var hashA uint64
	set := make([]reader.File, 0, len(dir1.Files))
	for _, a := range dir1.Files {
		if a.Info.IsDir() {
			continue
		}
		hashA = hashCalc(a.Name)
		for _, b := range intersec.Files {
			hashB := hashCalc(b.Name)
			if hashA != hashB && filepath.Base(a.Name) != filepath.Base(b.Name) {
				set = append(set, a)
				// log.Printf("File: %s hash: %d AND File: %s hash: %d - add to sync", a.Name, hashA, b.Name, hashB)
			} else {
				// log.Printf("File: %s hash: %d AND File: %s hash: %d - the same, not sync", a.Name, hashA, b.Name, hashB)
			}
		}
	}

	return reader.ListFiles{Files: set}
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

/*
Does not work...
The idea is that this func should work without finding intersection
*/
func Distinction2(dir1, dir2 reader.ListFiles) reader.ListFiles {
	if len(dir1.Files) == 0 || len(dir2.Files) == 0 {
		return reader.ListFiles{}
	}
	var hashA uint64
	set := make([]reader.File, 0, len(dir1.Files))
	for _, a := range dir1.Files {
		if a.Info.IsDir() {
			continue
		}
		hashA = hashCalc(a.Name)
		for _, b := range dir2.Files {
			if b.Info.IsDir() {
				continue
			}
			if hashA != hashCalc(b.Name) && filepath.Base(a.Name) != filepath.Base(b.Name) {
				set = append(set, a)
			}
		}
	}

	return reader.ListFiles{Files: set}
}
