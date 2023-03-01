package main

import (
	"context"
	"final_task/internal/reader"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

const (
	// source_folder      = "/workspaces/rebrain-go/final_task/testbed/src_folder/dir1/subdir1/"
	source_folder      = "/workspaces/rebrain-go/final_task/testbed/src_folder/"
	destination_folder = "/workspaces/rebrain-go/final_task/testbed/dst_folder/"
	// source_folder      = "./testbed/src_folder/"
	// destination_folder = "./testbed/dst_folder/"
)

func main() {
	var wg sync.WaitGroup
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)

	reader.LogSetup(6)
	wg.Add(1)
	go reader.Sync(ctx, &wg, source_folder, destination_folder, true)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	// go func() {
	for sig := range c {
		// sig is a ^C, handle it
		fmt.Printf(sig.String())
		cancelFunc()
		wg.Wait()
		return
	}
	// }()

}
