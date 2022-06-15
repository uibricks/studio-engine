package main

import (
	"github.com/uibricks/studio-engine/internal/app/project/di"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := di.InitializeApp()
	checkErr(err)

	app.Start(checkErr)

	<-interrupt()

	app.Shutdown()
}

func interrupt() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	return interrupt
}

func checkErr(err error) {
	if err != nil {
		logger.Log.Fatal("failed to start app", zap.Error(err))
	}
}
