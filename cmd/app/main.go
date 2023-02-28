package main

import (
	"final_task/internal/reader"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

const (
	// source_folder      = "/workspaces/rebrain-go/final_task/testbed/src_folder/dir1/subdir1/"
	source_folder      = "/workspaces/rebrain-go/final_task/testbed/src_folder/"
	destination_folder = "/workspaces/rebrain-go/final_task/testbed/dst_folder/"
	// source_folder      = "./testbed/src_folder/"
	// destination_folder = "./testbed/dst_folder/"
)

func main() {
	logger.SetLevel(6)
	entry := logrus.NewEntry(logger)

	var err error
	reader.LogSetup(6)
	dir1 := reader.NewWalker()
	dir1, err = dir1.WalkDir(source_folder)
	if err != nil {
		entry.Error(err)
	}
	// dir1.PrintData()

	dir2 := reader.NewWalker()
	dir2, err = dir2.WalkDir(destination_folder)
	if err != nil {
		entry.Error(err)
	}
	// dir2.PrintData()

	// distinc := reader.Distinction(*dir1, *dir2)
	// distinc2 := reader.Distinction(*dir2, *dir1)

	// fmt.Println("Dir1->Dir2: Distinction")
	// distinc.PrintData()
	// fmt.Println("Dir2->Dir1: Distinction")
	// distinc2.PrintData()

	reader.Sync(dir1, dir2)
	reader.Sync(dir2, dir1)
}
