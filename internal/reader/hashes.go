package reader

import (
	"bufio"
	"log"
	"os"

	"github.com/zeebo/xxh3"
)

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
