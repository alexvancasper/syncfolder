package reader

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func (l *ListFiles) WalkDir(path string) (*ListFiles, error) {
	var output []string

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			entry.Error(err)
			return err
		}
		output = append(output, path)
		return nil
	})

	result := ListFiles{Basepath: path, MapHash: make(map[uint64]File), MapName: make(map[string]File)}

	for i := 0; i < len(output); i++ {
		info, err := os.Stat(output[i])
		if err != nil {
			entry.Error(err)
			continue
		}
		if info.IsDir() {
			continue
		}
		hashVal := hashCalc(output[i])
		f := File{
			Name: output[i],
			Hash: hashVal,
			Info: info,
		}
		result.MapHash[hashVal] = f
		if name := MakeNameKey(path, output[i]); len(name) != 0 {
			result.MapName[name] = f
		}
	}
	return &result, nil
}

func (l *ListFiles) PrintData() {
	for key, val := range l.MapName {
		entry.Debugf("Key: %s -> %s, %d, %d, %d\n", key, val.Name, val.Hash, val.Info.Size(), val.Info.ModTime().UnixNano())
	}
}

func LogSetup(level logrus.Level) {
	logger.SetLevel(level)
	entry = logrus.NewEntry(logger)
}

func NewWalker() *ListFiles {
	return &ListFiles{Basepath: "", MapHash: make(map[uint64]File), MapName: make(map[string]File)}
}

func MakeNameKey(path, filepath string) string {
	res, _ := strings.CutPrefix(filepath, path)
	return res
}

// func InitLog() {
// 	// Set the output style, with only two style logrus.jsonformatter {} and logrus.textformatter {}
// 	log.SetFormatter(&log.TextFormatter{})
// 	log.SetOutput(os.Stdout)
// 		// Set Output, default is stderr, can be any I.Writer, such as file * Os.File
// 	 file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 	 writers := []io.Writer{
// 		file,
// 		os.Stdout}
// 		// Write the file and screen at the same time
// 	 fileAndStdoutWriter := io.MultiWriter(writers...)
// 	if err == nil {
// 	   log.SetOutput(fileAndStdoutWriter)
// 	} else {
// 	   log.Info("failed to log to file.")
// 	}
// 		// Set the lowest loglevel
// 	log.SetLevel(log.InfoLevel)
//  }
