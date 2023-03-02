package reader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	basePath          = "./testbed"
	sourceFolder      = "src_folder"
	destinationFolder = "dst_folder"
)

var (
	testPath = []string{"", "dir1"} // "./dir1"
)

/*
size - how many bytes need to generate. in bytes
*/
// func randData(size int) []byte {
// 	buf := make([]byte, size)
// 	_, err := rand.Read(buf)
// 	if err != nil {
// 		log.Printf("unable to create test dir tree: %v\n", err)
// 	}
// 	return buf
// }

func createFiles(root string, dirs []string, count int, mode os.FileMode) {
	var path string
	// source_folder := filepath.Join(basePath, sourceFolder) // ./testbad/src_folder/
	source_folder := root
	// info, _ := os.Stat(source_folder)
	for _, dir := range dirs {
		path = filepath.Join(source_folder, dir) // ./testbad/src_folder/''  , ./testbad/src_folder/dir1 ,  ./testbad/src_folder/dir1/subdir1 , ./testbad/src_folder/dir1/subdir1/subsubdir1
		source_folder = path
		for i := 0; i < count; i++ {
			fileName := fmt.Sprintf("filename_%d.txt", i+1)
			filePath := filepath.Join(path, fileName)
			err := os.WriteFile(filePath, []byte(filePath), mode)
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

// func removeFake(path string) error {
// 	os.Chmod(path, os.ModeDir|fs.FileMode(os.O_RDWR))
// 	return os.RemoveAll(path)
// }

// func createFake(root string, mode os.FileMode) error {
// 	path := filepath.Join(root, "fakeDir")
// 	err := os.MkdirAll(path, 0777)
// 	os.Chmod(path, 0777)
// 	if err != nil {
// 		return fmt.Errorf("error creating %s directory: %v\n", root, err)
// 	}
// 	createFiles(root, []string{"fakeDir"}, 1, 0666)
// 	os.Chmod(path, 0777)
// 	return nil
// }

func TestWalkDir(t *testing.T) {
	source_folder := filepath.Join(basePath, sourceFolder)
	createTestDirTree(source_folder, os.FileMode.Perm(0777))
	createFiles(source_folder, testPath, 1, os.FileMode.Perm(0666))

	req := require.New(t)
	sourceDir := NewWalker()
	sourceDir, err := sourceDir.WalkDir(source_folder)
	if err != nil {
		req.Fail("TestWalkDir", err)
	}
	req.Len(sourceDir.MapHash, 2)
	req.Len(sourceDir.MapName, 2)
	req.Equal(source_folder, sourceDir.Basepath)
	req.EqualValues(6514270587029032945, sourceDir.MapHash[6514270587029032945].Hash)

	err = removeAll(basePath)
	if err != nil {
		log.Printf("%s", err)
	}
}

func TestWalkDirWithError(t *testing.T) {
	source_folder := filepath.Join(basePath, sourceFolder)
	createTestDirTree(source_folder, os.FileMode.Perm(0777))
	createFiles(source_folder, testPath, 1, os.FileMode.Perm(0666))
	os.Chmod(filepath.Join(source_folder, "dir1"), os.FileMode.Perm(0600))

	req := require.New(t)
	sourceDir := NewWalker()
	sourceDir, err := sourceDir.WalkDir(source_folder)
	if err != nil {
		req.Failf("TestWalkDirWithError: %+v", err.Error())
	}
	req.Len(sourceDir.MapHash, 1)
	req.Equal(source_folder, sourceDir.Basepath)
	req.EqualValues(6514270587029032945, sourceDir.MapHash[6514270587029032945].Hash)

	os.Chmod(filepath.Join(source_folder, "dir1"), os.FileMode.Perm(0777))
	err = removeAll(basePath)
	if err != nil {
		log.Printf("%s", err)
	}
}

func TestMakeNameKey(t *testing.T) {
	req := require.New(t)
	source_folder := filepath.Join(basePath, sourceFolder)
	filepath := filepath.Join(source_folder, "dir1/filename_1.txt")
	req.Equal("/dir1/filename_1.txt", makeNameKey(source_folder, filepath))
}
