package integration

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type services struct {
	pg      testcontainers.Container
	quora   testcontainers.Container
	network testcontainers.Network
}

func spawnPg(ctx context.Context, network string) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Networks:     []string{network},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		Env: map[string]string{
			"POSTGRES_USER":     "pgquora",
			"POSTGRES_PASSWORD": "pgquora",
			"POSTGRES_DB":       "pgquora",
		},
	}

	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	var resolveErr = func(err error) {
		if err != nil {
			log.Fatalf("resolve error: %v", err)
		}
	}

	port, err := pgC.MappedPort(ctx, "5432/tcp")
	resolveErr(err)

	os.Setenv("PG_HOST", "localhost")
	os.Setenv("PG_PORT", port.Port())
	os.Setenv("PG_USER", "pgquora")
	os.Setenv("PG_PASSWORD", "pgquora")
	os.Setenv("PG_DBNAME", "pgquora")

	return pgC
}

func spawnQuora(ctx context.Context, pgC testcontainers.Container, network string) testcontainers.Container {
	ip, err := pgC.ContainerIP(ctx)
	if err != nil {
		log.Fatalf("failed get pgC ip: %v", err)
	}

	req := testcontainers.ContainerRequest{
		Image:        "quora:local",
		ExposedPorts: []string{"3000"},
		Networks:     []string{network},
		Env: map[string]string{
			"PG_HOST":            ip,
			"PG_PORT":            "5432",
			"PG_USER":            "pgquora",
			"PG_PASSWORD":        "pgquora",
			"PG_DBNAME":          "pgquora",
			"APP_PORT":           ":3000",
			"JWT_ACCESS_SECRET":  uuid.NewString(),
			"JWT_REFRESH_SECRET": uuid.NewString(),
		},
		WaitingFor: wait.ForLog("success run migrations"),
	}

	quora, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("failed spawn quora container: %v", err)
	}

	return quora
}

func spawnServices(ctx context.Context) (services, func()) {
	var (
		pg    testcontainers.Container
		quora testcontainers.Container
	)

	networkRequest := testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:       "test",
			Attachable: true,
			Internal:   false,
		},
	}

	network, err := testcontainers.GenericNetwork(ctx, networkRequest)
	if err != nil {
		log.Fatalf("failed create network: %v", err)
	}

	pg = spawnPg(ctx, networkRequest.Name)

	if pg.IsRunning() {
		quora = spawnQuora(ctx, pg, networkRequest.Name)
	}

	svc := services{
		pg:      pg,
		quora:   quora,
		network: network,
	}

	return svc, func() {
		if err := svc.quora.Terminate(ctx); err != nil {
			log.Fatalf("fail termintate quora: %v", err)
		}
		if err := svc.pg.Terminate(ctx); err != nil {
			log.Fatalf("fail terminate pg container: %v", err)
		}

		if err := svc.network.Remove(ctx); err != nil {
			log.Fatalf("fail remove network: %v", err)
		}
	}
}
