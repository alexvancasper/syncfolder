package main

import (
	"context"
	"final_task/internal/reader"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"

	"final_task/internal/config"

	"github.com/sirupsen/logrus"
)

var AppConfig *config.Config

// const (
// 	source_folder      = "/workspaces/rebrain-go/final_task/testbed/src_folder/"
// 	destination_folder = "/workspaces/rebrain-go/final_task/testbed/dst_folder/"
// 	twoWay             = true
// 	logFile            = "log.txt"
// 	configPath         = "./config/config.yml"
// )

func main() {
	AppConfig = config.ReadConfig(os.Args[1])

	logger := logrus.New()
	file, err := os.OpenFile(AppConfig.Options.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("log file open error %s: %v", AppConfig.Options.LogFile, err)
		return
	}
	logger.SetOutput(file)
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
	logger.SetLevel(logrus.TraceLevel)
	logger.SetLevel(logrus.Level(AppConfig.Options.Debug))
	logEntry := logrus.NewEntry(logger)
	logEntry.Infof("Start service %s", AppConfig.Service.Name)
	defer file.Close()

	var wg sync.WaitGroup
	ctx, cancelFunc := context.WithCancel(context.Background())

	wg.Add(1)
	go reader.Sync(ctx, &wg, logEntry, AppConfig)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	for sig := range c {
		logEntry.Info(fmt.Errorf("received signal %s", sig.String()))
		cancelFunc()
		wg.Wait()
		file.Close()
		break
	}
	return
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}
