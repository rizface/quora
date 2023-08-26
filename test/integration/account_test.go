package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
			"PG_HOST":     ip,
			"PG_PORT":     "5432",
			"PG_USER":     "pgquora",
			"PG_PASSWORD": "pgquora",
			"PG_DBNAME":   "pgquora",
			"APP_PORT":    ":3000",
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

func spawnServices(ctx context.Context) services {
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

	return services{
		pg:      pg,
		quora:   quora,
		network: network,
	}
}

func TestCreateAccount(t *testing.T) {
	var (
		ctx      = context.Background()
		services = spawnServices(ctx)
	)

	defer func() {
		if err := services.quora.Terminate(ctx); err != nil {
			log.Fatalf("fail termintate quora: %v", err)
		}
		if err := services.pg.Terminate(ctx); err != nil {
			log.Fatalf("fail terminate pg container: %v", err)
		}

		if err := services.network.Remove(ctx); err != nil {
			log.Fatalf("fail remove network: %v", err)
		}
	}()

	type scenario struct {
		name             string
		payload          map[string]interface{}
		code             int
		checkExpectation func(t *testing.T, resp *http.Response)
	}

	scenarios := []scenario{
		{
			name: "success create one employee",
			payload: map[string]interface{}{
				"username": "rizface",
				"password": "12345678",
				"email":    "nice@gmail.com",
			},
			code: http.StatusOK,
		},
		{
			name: "failed create one employee - empty required field",
			payload: map[string]interface{}{
				"username": "",
				"password": "",
				"email":    "",
			},
			code: http.StatusBadRequest,
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				type response struct {
					Data map[string]string `json:"data"`
				}

				var result response

				err := json.NewDecoder(resp.Body).Decode(&result)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, "cannot be blank", result.Data["username"])
				assert.Equal(t, "cannot be blank", result.Data["email"])
				assert.Equal(t, "cannot be blank", result.Data["password"])
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			client := &http.Client{}

			url, err := services.quora.Endpoint(ctx, "")
			if err != nil {
				t.Error(err)
			}

			payload, err := json.Marshal(s.payload)
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/%s", url, "accounts"), bytes.NewReader(payload))
			if err != nil {
				t.Error(err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
			}

			if s.checkExpectation != nil {
				s.checkExpectation(t, resp)
			}
		})
	}
}
