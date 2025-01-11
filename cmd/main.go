package main

import (
	"context"
	"go.uber.org/dig"
	"os"
	"os/signal"
	"shop-bot/di"
	"shop-bot/internal/bot"
	"shop-bot/utils"
	"sync"
	"syscall"
)

func startBot(container *dig.Container) {
	err := container.Invoke(bot.StartBot)

	utils.PanicIfError(err)
}

func main() {
	shutdownContext, cancel := context.WithCancel(context.Background())

	defer cancel()

	var wg sync.WaitGroup

	container := di.BuildContainer()

	container = di.AppendDependenciesToContainer(container, []di.Dependency{
		{
			Constructor: func() context.Context {
				return shutdownContext
			},
			Interface: nil,
			Token:     "ShutdownContext",
		},
		{
			Constructor: func() *sync.WaitGroup {
				return &wg
			},
			Interface: nil,
			Token:     "ShutdownWaitGroup",
		},
	})

	wg.Add(1)

	go startBot(container)
	
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh

	cancel()

	wg.Wait()
}
