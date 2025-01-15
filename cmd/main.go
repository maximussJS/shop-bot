package main

import (
	"context"
	"go.uber.org/dig"
	"os"
	"os/signal"
	"shop-bot/di"
	"shop-bot/internal/bot"
	"shop-bot/internal/http"
	"shop-bot/utils"
	"sync"
	"syscall"
)

type runAppDependencies struct {
	dig.In

	Bot        bot.IBot         `name:"Bot"`
	HttpServer http.IHttpServer `name:"HttpServer"`
}

func start(container *dig.Container) {
	err := container.Invoke(func(deps runAppDependencies) {
		go deps.Bot.Start()

		deps.HttpServer.Start()
	})

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

	go start(container)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh

	cancel()

	wg.Wait()
}
