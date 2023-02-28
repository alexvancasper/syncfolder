package reader

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	basePath          = "./testbad"
	sourceFolder      = "src_folder"
	destinationFolder = "dst_folder"

	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)

var (
	// testPath = []string{""}
	testPath = []string{"", "dir1", "subdir1", "subsubdir1"} // "./dir1/subdir1/subsubdir1"
)

/*
size - how many bytes need to generate. in bytes
*/
func randData(size int) []byte {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		log.Printf("unable to create test dir tree: %v\n", err)
	}
	return buf
}

func createFiles(root string, dirs []string, count int, mode os.FileMode) {
	var path string
	// source_folder := filepath.Join(basePath, sourceFolder) // ./testbad/src_folder/
	source_folder := root
	for _, dir := range dirs {
		path = filepath.Join(source_folder, dir) // ./testbad/src_folder/''  , ./testbad/src_folder/dir1 ,  ./testbad/src_folder/dir1/subdir1 , ./testbad/src_folder/dir1/subdir1/subsubdir1
		source_folder = path
		for i := 1; i <= count; i++ {
			fileName := fmt.Sprintf("filename_%d.txt", i)
			filePath := filepath.Join(path, fileName)
			err := os.WriteFile(filePath, []byte("h"), mode)
			if err != nil {
				log.Printf("Unable to write file: %v", err)
				continue
			}
		}
	}
}

func createTestDirTree(path string, mode os.FileMode) error {
	var dirs string
	for _, dir := range testPath {
		dirs = filepath.Join(dirs, dir)
	}
	err := os.MkdirAll(filepath.Join(path, dirs), mode)
	if err != nil {
		return fmt.Errorf("error creating %s directory: %v\n", path, err)
	}
	return nil
}

func removeAll(path string) error {
	return os.RemoveAll(path)
}

func removeFake(path string) error {
	os.Chmod(path, os.ModeDir|OS_USER_RWX)
	return os.RemoveAll(path)
}

func createFake(root string, mode os.FileMode) error {
	path := filepath.Join(root, "fakeDir")
	err := os.MkdirAll(path, os.ModeDir|OS_USER_RWX)
	if err != nil {
		return fmt.Errorf("error creating %s directory: %v\n", root, err)
	}
	createFiles(root, []string{"fakeDir"}, 10, mode)
	os.Chmod(path, mode)
	return nil
}

func equals(dir1, dir2 []File) bool {
	return true
}

func TestWalk(t *testing.T) {
	source_folder := filepath.Join(basePath, sourceFolder)
	createTestDirTree(source_folder, os.ModeDir|(OS_USER_RWX|OS_ALL_R))
	createFiles(source_folder, testPath, 2, os.ModeDir|(OS_USER_RWX|OS_ALL_R))
	createFake(filepath.Join(basePath, sourceFolder), os.ModeDir|(OS_USER_R|OS_GROUP_R))

	sourceDir := NewWalker()
	sourceDir, err := sourceDir.WalkDir(source_folder)
	if err != nil {
		log.Printf("%s", err)
	}
	sourceDir.PrintData()
	// var etalon ListFiles
	// copy(etalon.Files, sourceDir.Files) // TODO: define etalon separately.
	// equals(sourceDir.Files, etalon.Files)

	err = removeFake(filepath.Join(basePath, sourceFolder, "fakeDir"))
	if err != nil {
		log.Printf("%s", err)
	}
	err = removeAll(basePath)
	if err != nil {
		log.Printf("%s", err)
	}
}
