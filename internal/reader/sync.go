package reader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var BUFFERSIZE int64

func Distinction(dir1, dir2 ListFiles) ListFiles {
	if len(dir1.MapHash) == 0 && len(dir2.MapHash) == 0 {
		return ListFiles{Basepath: "", MapHash: make(map[uint64]File), MapName: make(map[string]File)}
	}

	set := ListFiles{Basepath: "", MapHash: make(map[uint64]File), MapName: make(map[string]File)}
	set.Basepath = dir1.Basepath
	for hashA, a := range dir1.MapHash {
		if b, ok := dir2.MapHash[hashA]; ok {
			//Hash found in A
			if MakeNameKey(dir1.Basepath, a.Name) == MakeNameKey(dir2.Basepath, b.Name) {
				//Name and place in the destination and source folder are the same. Means files completly identically
				entry.Infof("Files %s and %s identially - not sync", a.Name, b.Name)
			} else {
				//Names or place is different - need to sync them
				entry.Infof("Files %s and %s content is the same but place is different - add to sync", a.Name, b.Name)
				// entry.Printf("Files %s and %s content is the same but place is different - add to sync", MakeNameKey(dir1.Basepath, a.Name), MakeNameKey(dir2.Basepath, b.Name))
				set.MapName[MakeNameKey(dir1.Basepath, a.Name)] = a
			}
		} else {
			//hash not found in B
			if b, ok := dir2.MapName[MakeNameKey(dir1.Basepath, a.Name)]; ok {
				//Name and place in the destination and source folder are the same. Need to add to sync the last changed file
				if a.Info.ModTime().UnixNano() > b.Info.ModTime().UnixNano() {
					entry.Infof("Files %s and %s First file is last changed. %s > %s. Add to sync first file - A", MakeNameKey(dir1.Basepath, a.Name), MakeNameKey(dir2.Basepath, b.Name), a.Info.ModTime(), b.Info.ModTime())
					set.MapName[MakeNameKey(dir1.Basepath, a.Name)] = a
				} else if a.Info.ModTime().UnixNano() < b.Info.ModTime().UnixNano() {
					entry.Warningf("Files %s and %s Second file is last changed. %s < %s. Not added to sync due to sync A->B. last changed file is B", a.Name, b.Name, a.Info.ModTime(), b.Info.ModTime())
					// set.MapName[MakeNameKey(dir2.Basepath, b.Name)] = b
				} else {
					entry.Errorf("Files %s and %s Last change time is the same. %s == %s - CRITICAL ISSUE!", MakeNameKey(dir1.Basepath, a.Name), MakeNameKey(dir2.Basepath, b.Name), a.Info.ModTime(), b.Info.ModTime())
				}
			} else {
				//Names or place is different - need to sync them
				entry.Infof("File %s is new - add to sync", MakeNameKey(dir1.Basepath, a.Name))
				set.MapName[MakeNameKey(dir1.Basepath, a.Name)] = a
			}
		}
	}
	return set
}

func Sync(dir1, dir2 *ListFiles) {
	//Dir1 -> Dir2
	distinc := Distinction(*dir1, *dir2)
	if len(distinc.MapHash) == 0 {
		entry.Debugf("Nothing to sync %s -> %s", dir1.Basepath, dir2.Basepath)
	}
	for key, val := range distinc.MapName {
		sourceFolder := filepath.Join(dir1.Basepath, key)[:len(filepath.Join(dir1.Basepath, key))-len(filepath.Base(key))]
		destination := filepath.Join(dir2.Basepath, key)
		pathDst := destination[:len(destination)-len(filepath.Base(key))]
		_, errStat := os.Stat(pathDst) //check is there is all necessary folders in the destination folder
		if errStat != nil {
			info, e := os.Stat(sourceFolder)
			if e != nil {
				entry.Error(e)
				continue
			}
			errMkdir := os.MkdirAll(pathDst, info.Mode()) // if no, create them here
			if errMkdir != nil {
				entry.Error(e)
				continue
			}
		}

		err := copy(val.Name, destination, val.Info.Size())
		if err != nil {
			entry.Errorf("File copying failed: %q\n", err)
			continue
		}
		entry.Infof("Copy file %s size %d bytes - done", key, val.Info.Size())
	}
}

func copy(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	infoDst, err := os.Stat(dst)
	if err == nil {
		newFilename := fmt.Sprintf("%s_%d", dst, infoDst.ModTime().Nanosecond())
		os.Rename(dst, newFilename)
		defer os.Remove(newFilename) //TODO: need to check that copy is sucessfully finished before removing
		// return fmt.Errorf("File %s already exists.", dst)
	}

	// destination, err := os.Create(dst)
	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceFileStat.Mode())
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	return err
}
