package reader

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashCalc(t *testing.T) {
	req := require.New(t)
	err := os.WriteFile("testFile", []byte("testFiletestFiletestFiletestFie"), os.FileMode.Perm(0666))
	if err != nil {
		req.Fail("TestHashCalc. Fail to write file", err)
	}
	req.EqualValues(uint64(0xd0d58c9e1f6584c4), hashCalc("testFile"))
	os.Remove("testFile")
}

func TestHashCalcErrorStat(t *testing.T) {
	req := require.New(t)
	initLog()
	req.EqualValues(uint64(0), hashCalc("testFile"))
}

func TestHashCalcErrorOpen(t *testing.T) {
	req := require.New(t)
	initLog()
	err := os.WriteFile("testFile", []byte("testFiletestFiletestFiletestFie"), os.FileMode.Perm(0666))
	if err != nil {
		req.Fail("TestHashCalc. Fail to write file", err)
	}
	os.Chmod("testFile", os.FileMode.Perm(0000))
	req.EqualValues(uint64(0), hashCalc("testFile"))
	os.Chmod("testFile", os.FileMode.Perm(0666))
	os.Remove("testFile")
}

func BenchmarkHashCalc(b *testing.B) {
	err := os.WriteFile("testFile", []byte("testFiletestFiletestFiletestFie"), os.FileMode.Perm(0666))
	if err != nil {
		fmt.Errorf("BenchmarkHashCalc. Fail to write file: %v", err)
		return
	}
	for i := 0; i < b.N; i++ {
		hashCalc("testFile")
	}
	os.Remove("testFile")
}
