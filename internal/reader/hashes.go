package reader

import (
	"bufio"
	"fmt"
	"os"

	"github.com/zeebo/xxh3"
)

/*
hashCalc - calculate xxh3 hash based on content of file
https://cyan4973.github.io/xxHash/
*/
func hashCalc(path string) uint64 {
	fileSize, err := os.Stat(path)
	if err != nil {
		fmt.Errorf("hashCalc: stat %s, err %v", path, err)
	}

	buf := make([]byte, 0, fileSize.Size())

	readFile, err := os.Open(path)
	if err != nil {
		fmt.Errorf("hashCalc: open %s, err %v", path, err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	for fileScanner.Scan() {
		buf = append(buf, fileScanner.Bytes()...)
	}
	return xxh3.Hash(buf)
}
