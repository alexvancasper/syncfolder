package main

import (
	"final_task/internal/reader"
	"log"
)

const (
	// source_folder      = "/workspaces/rebrain-go/final_task/testbed/src_folder/dir1/subdir1/"
	source_folder      = "./testbed/src_folder/"
	destination_folder = "./testbed/dst_folder/"
)

func main() {
	var err error
	dir1 := reader.NewWalker()
	dir1, err = dir1.Walk(source_folder)
	if err != nil {
		log.Printf("%s", err)
	}
	dir1.PrintData()

	dir2 := reader.NewWalker()
	dir2, err = dir2.Walk(destination_folder)
	if err != nil {
		log.Fatalf("%s", err)
	}
	dir2.PrintData()

}
