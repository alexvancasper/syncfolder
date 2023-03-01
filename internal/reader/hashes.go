package reader

import (
	"bufio"
	"fmt"
	"os"

	"github.com/zeebo/xxh3"
)

func hashCalc(path string) uint64 {
	fileSize, err := os.Stat(path)
	if err != nil {
		entry.Error(fmt.Errorf("hashCalc: stat %s, err %v", path, err))
	}

	buf := make([]byte, 0, fileSize.Size())

	readFile, err := os.Open(path)
	if err != nil {
		entry.Error(fmt.Errorf("hashCalc: open %s, err %v", path, err))
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	for fileScanner.Scan() {
		buf = append(buf, fileScanner.Bytes()...)
	}
	return xxh3.Hash(buf)
}
