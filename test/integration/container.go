package integration

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type services struct {
	pg      testcontainers.Container
	quora   testcontainers.Container
	network testcontainers.Network
	jgr     testcontainers.Container
}

func resolveErr(err error) {
	if err != nil {
		log.Fatalf("resolve error: %v", err)
	}
}

func spawnPg(ctx context.Context, network string) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Networks:     []string{network},
		WaitingFor:   wait.ForListeningPort("5432"),
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

	port, err := pgC.MappedPort(ctx, "5432/tcp")
	resolveErr(err)

	os.Setenv("PG_HOST", "localhost")
	os.Setenv("PG_PORT", port.Port())
	os.Setenv("PG_USER", "pgquora")
	os.Setenv("PG_PASSWORD", "pgquora")
	os.Setenv("PG_DBNAME", "pgquora")

	return pgC
}

func spawnQuora(ctx context.Context, pgC, jgr testcontainers.Container, network string) testcontainers.Container {
	ip, err := pgC.ContainerIP(ctx)
	if err != nil {
		log.Fatalf("failed get pgC ip: %v", err)
	}

	jgrIp, err := jgr.ContainerIP(ctx)
	resolveErr(err)

	req := testcontainers.ContainerRequest{
		Image:        "quora:local",
		ExposedPorts: []string{"3000"},
		Networks:     []string{network},
		Env: map[string]string{
			"PG_HOST":             ip,
			"PG_PORT":             "5432",
			"PG_USER":             "pgquora",
			"PG_PASSWORD":         "pgquora",
			"PG_DBNAME":           "pgquora",
			"APP_PORT":            ":3000",
			"JWT_ACCESS_SECRET":   "access secret",
			"JWT_REFRESH_SECRET":  "refresh secret",
			"JAEGER_EXPORTER_URL": fmt.Sprintf("http://%s/api/traces", jgrIp),
		},
		WaitingFor: wait.ForListeningPort("3000"),
	}
	fmt.Println("JAEGER: ", fmt.Sprintf("http://%s:14268/api/traces", jgrIp))
	quora, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           log.Default(),
	})
	if err != nil {
		log.Fatalf("failed spawn quora container: %v", err)
	}

	return quora
}

func spawnJaeger(ctx context.Context, network string) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image: "jaegertracing/all-in-one:1.49",
		ExposedPorts: []string{
			"6831/udp",
			"6832/udp",
			"5778",
			"16686",
			"4317",
			"4318",
			"14250",
			"14268",
			"14269",
			"9411",
		},
		WaitingFor: wait.ForListeningPort("14268"),
		Networks:   []string{network},
	}

	jgr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	resolveErr(err)

	return jgr
}

func spawnServices(ctx context.Context) (services, func()) {
	var (
		pg    testcontainers.Container
		quora testcontainers.Container
		jgr   testcontainers.Container
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
	jgr = spawnJaeger(ctx, networkRequest.Name)

	if pg.IsRunning() && jgr.IsRunning() {
		quora = spawnQuora(ctx, pg, jgr, networkRequest.Name)
	}

	if !quora.IsRunning() {
		resolveErr(errors.New("quora is not running"))
	}

	svc := services{
		pg:      pg,
		quora:   quora,
		network: network,
		jgr:     jgr,
	}

	return svc, func() {
		if err := svc.quora.Terminate(ctx); err != nil {
			log.Fatalf("fail termintate quora: %v", err)
		}
		if err := svc.pg.Terminate(ctx); err != nil {
			log.Fatalf("fail terminate pg container: %v", err)
		}

		if err := svc.jgr.Terminate(ctx); err != nil {
			log.Fatalf("failed terminal jaeger container: %v", err)
		}

		if err := svc.network.Remove(ctx); err != nil {
			log.Fatalf("fail remove network: %v", err)
		}
	}
}
