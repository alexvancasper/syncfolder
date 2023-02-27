package main

import (
	"final_task/internal/reader"
	"final_task/internal/sync"
	"fmt"
	"log"
)

const (
	// source_folder      = "/workspaces/rebrain-go/final_task/testbed/src_folder/dir1/subdir1/"
	// source_folder      = "/workspaces/rebrain-go/final_task/testbed/src_folder/"
	source_folder      = "./testbed/src_folder/"
	destination_folder = "./testbed/dst_folder/"
)

func main() {
	var err error
	dir1 := reader.NewWalker()
	dir1, err = dir1.WalkDir(source_folder)
	if err != nil {
		log.Printf("%s", err)
	}
	// dir1.PrintData()

	dir2 := reader.NewWalker()
	dir2, err = dir2.WalkDir(destination_folder)
	if err != nil {
		log.Printf("%s", err)
	}
	// dir2.PrintData()

	distinc := sync.Distinction(*dir1, *dir2)
	distinc2 := sync.Distinction(*dir2, *dir1)

	fmt.Println("Dir1->Dir2: Distinction")
	distinc.PrintData()
	fmt.Println("Dir2->Dir1: Distinction")
	distinc2.PrintData()

}
