package reader

import (
	"context"
	"final_task/internal/config"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func initLog() {
	logger := logrus.New()
	logger.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat: "Mon Jan 2 15:04:05 MST 2006",
		FullTimestamp:   true,
		// DisableLevelTruncation: true, // log level field configuration
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	}
	logger.SetFormatter(formatter)
	logger.SetLevel(6)
	entry = logrus.NewEntry(logger)
}

func initDir() {
	os.MkdirAll("testbed2/src_folder/dir1/", os.FileMode.Perm(0777))
	os.MkdirAll("testbed2/dst_folder/dir2/", os.FileMode.Perm(0777))
	//Identical files
	os.WriteFile("testbed2/src_folder/file1", []byte("file1"), os.FileMode.Perm(0777))
	os.WriteFile("testbed2/dst_folder/file1", []byte("file1"), os.FileMode.Perm(0777))
	//Same content. Path is different
	os.WriteFile("testbed2/src_folder/file2", []byte("file2"), os.FileMode.Perm(0777))
	os.WriteFile("testbed2/dst_folder/file21", []byte("file2"), os.FileMode.Perm(0777))
	//Different content, Path is the same, dst is last changed
	os.WriteFile("testbed2/src_folder/file3", []byte("file3"), os.FileMode.Perm(0777))
	time.Sleep(1 * time.Second)
	os.WriteFile("testbed2/dst_folder/file3", []byte("file31"), os.FileMode.Perm(0777))
	//Different content, Path is the same, src is last changed
	os.WriteFile("testbed2/dst_folder/file4", []byte("file4"), os.FileMode.Perm(0777))
	time.Sleep(1 * time.Second)
	os.WriteFile("testbed2/src_folder/file4", []byte("file41"), os.FileMode.Perm(0777))
	//The same time of changed
	os.WriteFile("testbed2/dst_folder/file5", []byte("file5"), os.FileMode.Perm(0777))
	os.WriteFile("testbed2/src_folder/file5", []byte("file51"), os.FileMode.Perm(0777))
	//New files
	os.WriteFile("testbed2/src_folder/file61", []byte("file61"), os.FileMode.Perm(0777))
	os.WriteFile("testbed2/dst_folder/file62", []byte("file62"), os.FileMode.Perm(0777))
}

func TestDistinction(t *testing.T) {
	initLog()
	initDir()
	defer os.RemoveAll("testbed2/")
	req := require.New(t)
	dir1 := NewWalker()
	dir1, err := dir1.WalkDir("testbed2/src_folder")
	if err != nil {
		req.Fail("TestDistinction: dir1", err)
	}
	dir2 := NewWalker()
	dir2, err = dir2.WalkDir("testbed2/dst_folder")
	if err != nil {
		req.Fail("TestDistinctionL dir2", err)
	}

	result := distinction(*dir1, *dir2)
	req.Equal(0, len(result.MapHash))
	_, file2 := result.MapName["/file2"]
	req.Equal(true, file2)
	_, file4 := result.MapName["/file4"]
	req.Equal(true, file4)
	_, file61 := result.MapName["/file61"]
	req.Equal(true, file61)
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}

func TestCopyNormal(t *testing.T) {
	os.MkdirAll("testbed2/src_folder/dir1/", os.FileMode.Perm(0777))
	os.MkdirAll("testbed2/dst_folder/dir2/", os.FileMode.Perm(0777))
	data := strings.Repeat("TEST\n", 4000)
	os.WriteFile("testbed2/src_folder/file1", []byte(data), os.FileMode.Perm(0777))
	defer os.RemoveAll("testbed2/")
	req := require.New(t)
	info, _ := os.Stat("testbed2/src_folder/file1")
	copied, err := copyFile("testbed2/src_folder/file1", "testbed2/dst_folder/file1", info.Size())
	req.Equal(nil, err)
	req.Equal(int64(len(data)), copied)

}

func TestCopySourceStat(t *testing.T) {
	req := require.New(t)
	copied, err := copyFile("testbed2/src_folder/file1", "testbed2/dst_folder/file1", 100)
	req.ErrorContains(err, "copy: stat source file:")
	req.Equal(int64(0), copied)
	os.RemoveAll("testbed2/")
}

func TestCopySourceSymlink(t *testing.T) {
	req := require.New(t)
	os.MkdirAll("testbed2/src_folder/dir1/", os.FileMode.Perm(0777))
	os.WriteFile("testbed2/src_folder/file1", []byte("123"), os.FileMode.Perm(0777))
	defer os.RemoveAll("testbed2/")

	target := "file1"
	symlink := "testbed2/src_folder/symFile1"
	os.Symlink(target, symlink)
	copied, err := copyFile(symlink, "testbed2/dst_folder/file1", 100)
	req.ErrorContains(err, "is not a regular file.")
	req.Equal(int64(0), copied)
}

func TestCopyDestinationCreate(t *testing.T) {
	req := require.New(t)
	os.MkdirAll("testbed2/src_folder/dir1/", os.FileMode.Perm(0777))
	os.WriteFile("testbed2/src_folder/file1", []byte("123"), os.FileMode.Perm(0777))
	defer os.RemoveAll("testbed2/")

	copied, err := copyFile("testbed2/src_folder/file1", "testbed2/dst_folder/file1", 100)
	req.ErrorContains(err, "copy: destination file create")
	req.Equal(int64(0), copied)
}

func TestCopyDestinationDuplicat(t *testing.T) {
	req := require.New(t)
	os.MkdirAll("testbed2/src_folder/dir1/", os.FileMode.Perm(0777))
	os.MkdirAll("testbed2/dst_folder/dir1/", os.FileMode.Perm(0777))
	os.WriteFile("testbed2/src_folder/file1", []byte("123"), os.FileMode.Perm(0777))
	os.WriteFile("testbed2/dst_folder/file1", []byte("123"), os.FileMode.Perm(0777))
	defer os.RemoveAll("testbed2/")

	copied, err := copyFile("testbed2/src_folder/file1", "testbed2/dst_folder/file1", 100)
	req.Equal(int64(3), copied)
	req.Equal(nil, err)
}

func TestSyncTwoFolders(t *testing.T) {
	req := require.New(t)
	initLog()
	os.MkdirAll("testbed2/src_folder/dir1", os.FileMode.Perm(0777))
	os.MkdirAll("testbed2/dst_folder/", os.FileMode.Perm(0777))
	os.WriteFile("testbed2/src_folder/file1", []byte("123"), os.FileMode.Perm(0777))
	os.WriteFile("testbed2/src_folder/dir1/dir1-file1", []byte("123"), os.FileMode.Perm(0777))
	defer os.RemoveAll("testbed2/")

	dir1 := NewWalker()
	dir1, err := dir1.WalkDir("testbed2/src_folder")
	if err != nil {
		req.Fail("TestDistinction: dir1", err)
	}
	dir2 := NewWalker()
	dir2, err = dir2.WalkDir("testbed2/dst_folder")
	if err != nil {
		req.Fail("TestDistinctionL dir2", err)
	}
	syncTwoFolders(dir1, dir2)
	_, errDstFile1 := os.Stat("testbed2/dst_folder/file1")
	_, errDstDir1File1 := os.Stat("testbed2/src_folder/dir1/dir1-file1")
	req.Equal(nil, errDstFile1)
	req.Equal(nil, errDstDir1File1)

}

func TestSync(t *testing.T) {
	req := require.New(t)
	os.MkdirAll("testbed2/src_folder/dir1", os.FileMode.Perm(0777))
	os.MkdirAll("testbed2/dst_folder/", os.FileMode.Perm(0777))
	os.MkdirAll("testbed2/logs/", os.FileMode.Perm(0777))
	os.WriteFile("testbed2/src_folder/file1", []byte("123"), os.FileMode.Perm(0777))
	os.WriteFile("testbed2/src_folder/dir1/dir1-file1", []byte("1234"), os.FileMode.Perm(0777))
	defer os.RemoveAll("testbed2/")

	var AppConfig config.Config
	AppConfig.Service.Name = "Test Service"
	AppConfig.Folders.SrcFolder = "testbed2/src_folder/"
	AppConfig.Folders.DstFolder = "testbed2/dst_folder/"
	AppConfig.Options.Debug = 6
	AppConfig.Options.Internal = 1
	AppConfig.Options.LogFile = "testbed2/logs/log.txt"
	AppConfig.Options.TwoWay = true

	logger := logrus.New()
	file, err := os.OpenFile(AppConfig.Options.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("log file open error %s: %v", AppConfig.Options.LogFile, err)
		return
	}
	logger.SetOutput(file)
	logger.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat: "Mon Jan 2 15:04:05 MST 2006",
		FullTimestamp:   true,
		// DisableLevelTruncation: true, // log level field configuration
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	}
	logger.SetFormatter(formatter)
	logger.SetLevel(logrus.Level(AppConfig.Options.Debug))
	logEntry := logrus.NewEntry(logger)
	logEntry.Infof("Start service %s", AppConfig.Service.Name)
	defer file.Close()

	var wg sync.WaitGroup
	ctx, cancelFunc := context.WithCancel(context.Background())

	wg.Add(1)
	go Sync(ctx, &wg, logEntry, &AppConfig)

	time.Sleep(3000 * time.Millisecond)
	_, errFile1 := os.Stat("testbed2/dst_folder/file1")
	req.Equal(nil, errFile1)
	_, errDirFile1 := os.Stat("testbed2/src_folder/dir1/dir1-file1")
	req.Equal(nil, errDirFile1)

	os.WriteFile("testbed2/src_folder/testFile", []byte("12345"), os.FileMode.Perm(0777))
	time.Sleep(2000 * time.Millisecond)
	_, errTestFile := os.Stat("testbed2/dst_folder/testFile")
	req.Equal(nil, errTestFile)

	os.WriteFile("testbed2/dst_folder/dstFile", []byte("123456"), os.FileMode.Perm(0777))
	time.Sleep(3000 * time.Millisecond)
	_, errDstFile := os.Stat("testbed2/src_folder/dstFile")
	req.Equal(nil, errDstFile)

	cancelFunc()
	wg.Wait()

	_, errLog := os.Stat("testbed2/logs/log.txt")
	req.Equal(nil, errLog)

}

func BenchmarkCopyFile(b *testing.B) {
	os.MkdirAll("testbed2/src_folder/dir1/", os.FileMode.Perm(0777))
	os.MkdirAll("testbed2/dst_folder/dir2/", os.FileMode.Perm(0777))
	data := strings.Repeat("TB\n", 100000000)
	os.WriteFile("testbed2/src_folder/file1", []byte(data), os.FileMode.Perm(0777))
	defer os.RemoveAll("testbed2/")

	fileInfo, _ := os.Stat("testbed2/src_folder/file1")
	var bufSize int64
	if fileInfo.Size() > int64(100e6) {
		bufSize = 100e6
	} else {
		bufSize = fileInfo.Size()
	}

	for i := 0; i < b.N; i++ {
		copyFile("testbed2/src_folder/file1", "testbed2/dst_folder/file1", bufSize)
	}

}
