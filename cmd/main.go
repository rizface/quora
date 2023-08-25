package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rizface/quora/account"
	"github.com/rizface/quora/provider"
)

func main() {
	var (
		ctx          = context.Background()
		dependencies = InitDependencies()
		app          = NewApp(dependencies)
	)

	listener := make(chan os.Signal, 1)
	signal.Notify(listener, os.Interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		for range listener {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			go func() {
				<-ctx.Done()

				fmt.Println("system is forced to shutdown due to timeout")

				os.Exit(0)
			}()

			if err := app.Stop(ctx); err != nil {
				fmt.Println("failed to stop the app gracefully")

				cancel()

				os.Exit(1)
			}

			cancel()

			fmt.Println("successfully stopped the app")

			os.Exit(0)
		}
	}()

	// start the app
	log.Fatal(app.Start())
}

type App struct {
	Deps    *Dependencies
	Account *account.Feature
}

func NewApp(d *Dependencies) *App {
	return &App{
		Deps:    d,
		Account: account.NewFeature(d.router),
	}
}

func (a *App) Start() error {
	a.Account.RegisterRoutes()

	err := a.Deps.server.ListenAndServe()

	return err
}

func (s *App) Stop(ctx context.Context) error {
	err := s.Deps.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

type Dependencies struct {
	server *http.Server
	router *chi.Mux
}

func InitDependencies() *Dependencies {
	router := provider.ProvideRouter()
	server := provider.ProviderServer(router)

	return &Dependencies{
		router: router,
		server: server,
	}
}
