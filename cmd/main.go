package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rizface/quora/account"
	"github.com/rizface/quora/provider"
	"github.com/rizface/quora/question"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
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
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
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

	runMigrations(dependencies.sql)

	// start the app
	log.Fatal(app.Start())
}

type App struct {
	Deps     *Dependencies
	Account  *account.Feature
	Question *question.Feature
}

func NewApp(d *Dependencies) *App {
	return &App{
		Deps:     d,
		Account:  account.NewFeature(d.router, d.sql, d.tracer),
		Question: question.NewFeature(d.router, d.sql),
	}
}

func (a *App) Start() error {
	a.Account.RegisterRoutes()
	a.Question.RegisterRoutes()

	err := a.Deps.server.ListenAndServe()

	return err
}

func (s *App) Stop(ctx context.Context) error {
	err := s.Deps.sql.Close()
	if err != nil {
		return err
	}
	log.Println("SQL connection closed")

	err = s.Deps.traceProvider.Shutdown(ctx)
	if err != nil {
		return err
	}
	log.Println("OTel TraceProvider shutdown")

	err = s.Deps.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

type Dependencies struct {
	server        *http.Server
	router        *chi.Mux
	sql           *sql.DB
	tracer        trace.Tracer
	traceProvider *sdktrace.TracerProvider
}

func InitDependencies() *Dependencies {
	router := provider.ProvideRouter()
	server := provider.ProviderServer(router)

	sql, err := provider.ProvideSQL()
	if err != nil {
		log.Fatal(err)
	}

	traceProvider, tracer, err := provider.ProvideOtel()
	if err != nil {
		log.Fatal(err)
	}

	return &Dependencies{
		router:        router,
		server:        server,
		sql:           sql,
		tracer:        tracer,
		traceProvider: traceProvider,
	}
}

func runMigrations(sql *sql.DB) {
	driver, err := postgres.WithInstance(sql, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed create pg instance: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", os.Getenv("PG_DBNAME"), driver)
	if err != nil {
		log.Fatalf("failed to create migration instance: %v", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed run migrations: %v", err)
	}

	log.Print("success run migrations")
}
