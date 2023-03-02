package reader

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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
	req.Equal(3, len(result.MapName))
	_, file2 := result.MapName["/file2"]
	req.Equal(true, file2)
	_, file4 := result.MapName["/file4"]
	req.Equal(true, file4)
	_, file61 := result.MapName["/file61"]
	req.Equal(true, file61)
	os.RemoveAll("testbed2/")
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}
