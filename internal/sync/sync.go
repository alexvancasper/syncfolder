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
				// log.Printf("File: %s hash: %d AND File: %s hash: %d - fully matched", a.Name, hashA, b.Name, hashB)
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
	intersec := Intersection(dir1, dir2) //TODO: this is does not work in case if there is no matched files
	if len(intersec.Files) == 0 {
		intersec = dir2
	}
	var hashA uint64
	set := make([]reader.File, 0, len(dir1.Files))
	for _, a := range dir1.Files {
		if a.Info.IsDir() {
			continue
		}
		hashA = hashCalc(a.Name)
		for _, b := range intersec.Files {
			hashB := hashCalc(b.Name)
			if filepath.Base(a.Name) != filepath.Base(b.Name) {
				set = append(set, a)
				log.Printf("File: %s hash: %d AND File: %s hash: %d - add to sync because file is new.", a.Name, hashA, b.Name, hashB)
			} else if hashA != hashB {
				set = append(set, a)
				log.Printf("File: %s hash: %d AND File: %s hash: %d - add to sync due to hashes are different", a.Name, hashA, b.Name, hashB)
			} else {
				log.Printf("File: %s hash: %d AND File: %s hash: %d - the same, not sync", a.Name, hashA, b.Name, hashB)
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
func Distinction3(dir1, dir2 reader.ListFiles) reader.ListFiles {
	if len(dir1.Files) == 0 || len(dir2.Files) == 0 {
		return reader.ListFiles{}
	}
	uniqHash := make(map[uint64]bool)
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
			if filepath.Base(a.Name) != filepath.Base(b.Name) {
				if _, ok := uniqHash[hashA]; !ok {
					set = append(set, a)
					log.Printf("File: %s hash: %d AND File: %s hash: %d - add to sync because file is new.", a.Name, hashA, b.Name, hashB)
					uniqHash[hashA] = true
				}
			} else if hashA != hashB {
				log.Printf("File: %s hash: %d AND File: %s hash: %d - hashes are different. Check what file is new?", a.Name, hashA, b.Name, hashB)
				//Need to find out the last changed file
				if a.Info.ModTime().UnixNano() > b.Info.ModTime().UnixNano() { // file A is last changed
					if _, ok := uniqHash[hashA]; !ok {
						set = append(set, a)
						log.Printf("File: %s hash: %d - is last changed", a.Name, hashA)
						uniqHash[hashA] = true
					}
				} else {
					if _, ok := uniqHash[hashB]; !ok {
						set = append(set, b)
						log.Printf("File: %s hash: %d - is last changed", b.Name, hashB)
						uniqHash[hashB] = true
					}
				}
			} else {
				log.Printf("File: %s hash: %d AND File: %s hash: %d - the same, not sync", a.Name, hashA, b.Name, hashB)
			}
		}
	}
	return reader.ListFiles{Files: set}
}

func Distinction2(dir1, dir2 reader.ListFiles) reader.ListFiles {
	if len(dir1.Files) == 0 || len(dir2.Files) == 0 {
		return reader.ListFiles{}
	}
	flag := false
	set := make([]reader.File, 0, len(dir1.Files))
	mapFiles := make(map[uint64]reader.File)
	for _, a := range dir1.Files {
		if a.Info.IsDir() {
			continue
		}
		log.Printf("File: %s - %d", a.Name, a.Hash)
		if b, ok := dir2.Map[a.Hash]; ok { //The same hash already exisst. means the content of the file is already there.
			if b.Info.Name() != a.Info.Name() {
				if a.Info.ModTime().UnixNano() > b.Info.ModTime().UnixNano() { // file A is last changed
					set = append(set, a)
					mapFiles[a.Hash] = a
					log.Printf("File: %s hash: %d - A is last changed", a.Name, a.Hash)
					// } else if a.Info.ModTime().UnixNano() < b.Info.ModTime().UnixNano() { // file B is last changed
					// 	set = append(set, b)
					// 	mapFiles[b.Hash] = b
					// 	log.Printf("File: %s hash: %d - B is last changed", a.Name, a.Hash)
				} else {
					// log.Printf("File: %s hash: %d AND File: %s hash: %d and last modification time is the same - not sync", a.Name, a.Hash, b.Name, b.Hash)
				}
			} else {
				log.Printf("File: %s hash: %d AND File: %s hash: %d - the same, not sync", a.Name, a.Hash, b.Name, b.Hash)
			}
		} else {
			// Different hash - different content
			flag = false
			for _, f := range dir2.Map {
				if f.Info.Name() == a.Info.Name() {
					//File names the same
					if a.Info.ModTime().UnixNano() > f.Info.ModTime().UnixNano() { // file A is last changed
						set = append(set, a)
						mapFiles[a.Hash] = a
						log.Printf("File: %s hash: %d - A is last changed", a.Name, a.Hash)
					} else if a.Info.ModTime().UnixNano() < f.Info.ModTime().UnixNano() { // file B is last changed
						// set = append(set, f)
						// mapFiles[hash] = f
						flag = true
						log.Printf("File: %s hash: %d - B is last changed", f.Name, f.Hash)
					} else {
						// log.Printf("File: %s hash: %d AND File: %s hash: %d and last modification time is the same - not sync", a.Name, a.Hash, f.Name, f.Hash)
					}
					break
				}
			}
			//Means files are differents and need to sync them
			if !flag {
				if _, ok := mapFiles[a.Hash]; !ok {
					set = append(set, a)
					mapFiles[a.Hash] = a
					log.Printf("File: %s hash: %d - add to sync", a.Name, a.Hash)
				}

			}
		}
	}
	return reader.ListFiles{Files: set, Map: mapFiles}
}
