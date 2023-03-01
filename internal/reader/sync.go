package reader

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"final_task/internal/config"

	"github.com/sirupsen/logrus"
)

var entry *logrus.Entry

/*
distinction - function to find non-common files between two folders (ListFiles).
INPUT:
dir1 - source folder  or "a"
dir2 - destination folder or "b"
OUTPUT:
set - set of non-intersection files.
*/
func distinction(dir1, dir2 ListFiles) ListFiles {
	set := ListFiles{Basepath: "", MapHash: make(map[uint64]File), MapName: make(map[string]File)}
	if len(dir1.MapHash) == 0 && len(dir2.MapHash) == 0 {
		return set
	}
	set.Basepath = dir1.Basepath
	for hashA, a := range dir1.MapHash {
		if b, ok := dir2.MapHash[hashA]; ok {
			//hash found
			if makeNameKey(dir1.Basepath, a.Name) == makeNameKey(dir2.Basepath, b.Name) {
				//Name and path in the destination and source folder are the same. Hence files completly identically
				entry.Debugf("Both files identical. Files SRC:%s and DST:%s - not sync", a.Name, b.Name)
			} else {
				//Names or path is different - need to sync them
				entry.Infof("The same content. Path is different. Files SRC:%s and DST:%s  - SRC add to sync", a.Name, b.Name)
				set.MapName[makeNameKey(dir1.Basepath, a.Name)] = a
			}
		} else {
			//hash is not found
			if b, ok := dir2.MapName[makeNameKey(dir1.Basepath, a.Name)]; ok {
				//Name and path in the destination and source folder are the same. Need to add to sync the last changed file
				if a.Info.ModTime().UnixNano() > b.Info.ModTime().UnixNano() {
					entry.Infof("Source file is last changed. Files SRC:%s and DST:%s DateTime of changing: %s > %s - SRC add to sync ", makeNameKey(dir1.Basepath, a.Name), makeNameKey(dir2.Basepath, b.Name), a.Info.ModTime(), b.Info.ModTime())
					set.MapName[makeNameKey(dir1.Basepath, a.Name)] = a
				} else if a.Info.ModTime().UnixNano() < b.Info.ModTime().UnixNano() {
					entry.Debugf("Destination file is last changed. Files SRC:%s and DST:%s DateTime of changing: %s < %s. Not added to sync due to sync SRC->DST ", a.Name, b.Name, a.Info.ModTime(), b.Info.ModTime())
				} else {
					entry.Errorf("Last change time is the same. Files SRC:%s and DST:%s DateTime of changing: %s == %s - CRITICAL ISSUE!", makeNameKey(dir1.Basepath, a.Name), makeNameKey(dir2.Basepath, b.Name), a.Info.ModTime(), b.Info.ModTime())
				}
			} else {
				//Names or place is different - need to sync them
				entry.Infof("Found new file %s - add to sync", makeNameKey(dir1.Basepath, a.Name))
				set.MapName[makeNameKey(dir1.Basepath, a.Name)] = a
			}
		}
	}
	return set
}

/*
Sync - starting function, starts work every interval.
*/
func Sync(ctx context.Context, wg *sync.WaitGroup, entryF *logrus.Entry, AppConfig *config.Config) {
	entry = entryF //Assign to global var for logging
	ticker := time.NewTicker(time.Duration(AppConfig.Options.Internal) * time.Minute)

	entry.Infoln("Sync start")
	entry.Debugf("Source: %s Destination: %s TwoWay: %t", AppConfig.Folders.Srcfolder, AppConfig.Folders.DstFolder, AppConfig.Options.TwoWay)
	for {
		time.Sleep(50 * time.Millisecond) //necessary to avoid high CPU usage, but in general it is not needed.
		select {
		case <-ctx.Done():
			ticker.Stop()
			entry.Infoln("Sync shutdown")
			wg.Done()
			return

		case <-ticker.C:
			entry.Debugf("Start to check folders")
			var err error
			dir1 := NewWalker()
			dir1, err = dir1.WalkDir(AppConfig.Folders.Srcfolder)
			if err != nil {
				entry.Error(fmt.Errorf("Sync.WalkDir source: %v", err))
			}
			entry.Infof("Number of files in source folder: %d", len(dir1.MapHash))
			dir2 := NewWalker()
			dir2, err = dir2.WalkDir(AppConfig.Folders.DstFolder)
			if err != nil {
				entry.Error(fmt.Errorf("Sync.WalkDir destination: %v", err))
			}
			entry.Infof("Number of files in destination folder: %d", len(dir1.MapHash))
			if AppConfig.Options.TwoWay {
				syncTwoFolders(dir1, dir2)
				syncTwoFolders(dir2, dir1)
			} else {
				syncTwoFolders(dir1, dir2)
			}
		}
	}

}

/*
syncTwoFolders - If need to copy files.
Creates folder with the source permissions in the dst folder in case if it needed.
Then directories created copy files.
if the folder is empty then it won't sync.
*/
func syncTwoFolders(dir1, dir2 *ListFiles) {
	distinc := distinction(*dir1, *dir2)
	if len(distinc.MapHash) == 0 {
		entry.Infof("Nothing to sync %s -> %s", dir1.Basepath, dir2.Basepath)
	}
	for fileName, fileInfo := range distinc.MapName {
		/*
			source file path: /src-folder/dir1/dir2/file1
			need to copy it to: /dst-folder, in this case need to create folders dir1/dir2.
			below code removes /src-folder/ and file1 from source file path and add destination path /dst-folder/dir1/dir2/
		*/
		sourceFolder := filepath.Join(dir1.Basepath, fileName)[:len(filepath.Join(dir1.Basepath, fileName))-len(filepath.Base(fileName))]
		destination := filepath.Join(dir2.Basepath, fileName)
		pathDst := destination[:len(destination)-len(filepath.Base(fileName))]
		_, errStat := os.Stat(pathDst) //check is there is all necessary folders in the destination folder
		if errStat != nil {
			info, e := os.Stat(sourceFolder)
			if e != nil {
				entry.Error(fmt.Errorf("stat source folder: %v", e))
				continue
			}
			errMkdir := os.MkdirAll(pathDst, info.Mode()) // if no, create them here
			if errMkdir != nil {
				entry.Error(fmt.Errorf("mkdir in destination folder: %v", e))
				continue
			}
		}

		err := copy(fileInfo.Name, destination, fileInfo.Info.Size())
		if err != nil {
			entry.Error(fmt.Errorf("file copying failed: %v", err))
			continue
		}
		entry.Infof("Copy file SRC:%s to DST:%s size %d bytes - done", fileInfo.Name, destination, fileInfo.Info.Size())
	}
}

/*
Simple copy files with predefined buffer.
buf should not be huge value.
Possible issue: if the file is huge then it takes whole RAM.
*/
func copy(src, dst string, bufSize int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("copy: stat source file: %v", err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("copy: open source file: %v", err)
	}
	defer source.Close()

	//If destination file is exist then it renames the old file. After copy is done the file will be removed.
	newFilename := ""
	infoDst, err := os.Stat(dst)
	if err == nil {
		newFilename = fmt.Sprintf("%s_%d", dst, infoDst.ModTime().Nanosecond())
	}

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceFileStat.Mode())
	if err != nil {
		return fmt.Errorf("copy: destination file create: %v", err)
	}
	defer destination.Close()

	buf := make([]byte, bufSize)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("copy: read source file %v", err)
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return fmt.Errorf("copy: write destination file: %v", err)
		}
	}

	if err == nil && len(newFilename) > 3 {
		os.Remove(newFilename)
	}

	return err
}
