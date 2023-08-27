package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/rizface/quora/provider"
	"github.com/stretchr/testify/assert"
)

type scenario struct {
	name             string
	payload          map[string]interface{}
	checkExpectation func(t *testing.T, resp *http.Response)
}

func TestCreateAccount(t *testing.T) {
	var (
		ctx      = context.Background()
		services = spawnServices(ctx)
	)

	db, err := provider.ProvideSQL()
	if err != nil {
		log.Fatalf("failed open connection to pg: %v", err)
	}

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

	scenarios := []scenario{
		{
			name: "success create one employee",
			payload: map[string]interface{}{
				"username": "rizface",
				"password": "12345678",
				"email":    "nice@gmail.com",
			},
			checkExpectation: func(t *testing.T, resp *http.Response) {
				var counter int

				err := db.
					QueryRowContext(ctx, `
					SELECT COUNT(id) as count from accounts WHERE username = $1
				`, "rizface").
					Scan(&counter)
				assert.Nil(t, err)
				assert.Equal(t, int(1), counter)
			},
		},
		{
			name: "failed create one employee - empty required field",
			payload: map[string]interface{}{
				"username": "",
				"password": "",
				"email":    "",
			},
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

func TestLogin(t *testing.T) {
	var (
		ctx      = context.Background()
		services = spawnServices(ctx)
		db, err  = provider.ProvideSQL()
	)

	if err != nil {
		log.Fatal(err)
	}

	ImportSQL(db, "../../testdata/account/login.sql")

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

	scenarios := []scenario{
		{
			name: "success login",
			payload: map[string]interface{}{
				"email":    "testlogin@gmail.com",
				"password": "testdata",
			},
			checkExpectation: func(t *testing.T, resp *http.Response) {
				type Data struct {
					Data map[string]interface{} `json:"data"`
				}

				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var data Data
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					log.Fatal(err)
				}

				assert.Equal(t, data.Data["email"], "testlogin@gmail.com")
				assert.Equal(t, data.Data["username"], "testlogin")
				assert.Equal(t, data.Data["id"], "f028ac5a-e4c9-442f-bf9a-86c024a79baa")
				assert.Equal(t, len(data.Data["tokens"].([]interface{})), 2)
			},
		},
		{
			name: "failed login - account not found",
			payload: map[string]interface{}{
				"email":    "notfound@gmail.com",
				"password": "testdata",
			},
			checkExpectation: func(t *testing.T, resp *http.Response) {
				type Data struct {
					Info string `json:"info"`
				}

				assert.Equal(t, http.StatusNotFound, resp.StatusCode)

				var data Data
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					log.Fatal(err)
				}

				assert.Equal(t, data.Info, "account not found")
			},
		},
		{
			name: "failed login - wrong credential",
			payload: map[string]interface{}{
				"email":    "testlogin@gmail.com",
				"password": "wrongpassword",
			},
			checkExpectation: func(t *testing.T, resp *http.Response) {
				type Data struct {
					Info string `json:"info"`
				}

				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

				var data Data
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					log.Fatal(err)
				}

				assert.Equal(t, data.Info, "wrong email / password")
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

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/%s", url, "accounts/login"), bytes.NewReader(payload))
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
