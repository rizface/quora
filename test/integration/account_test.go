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
)

func TestCreateAccount(t *testing.T) {
	var (
		ctx      = context.Background()
		services = spawnServices(ctx)
	)

	// db, err := provider.ProvideSQL()
	// if err != nil {
	// 	log.Fatalf("failed open connection to pg: %v", err)
	// }

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
			// checkExpectation: func(t *testing.T, resp *http.Response) {
			// 	var counter int

			// 	err := db.
			// 		QueryRowContext(ctx, `
			// 		SELECT COUNT(id) as count from accounts WHERE username = $1
			// 	`, "rizface").
			// 		Scan(&counter)
			// 	assert.Nil(t, err)
			// 	assert.Equal(t, int(1), counter)
			// },
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
