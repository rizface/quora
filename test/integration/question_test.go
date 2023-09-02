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

func TestCreateNewQuestion(t *testing.T) {
	var (
		ctx               = context.Background()
		services, cleaner = spawnServices(ctx)
	)

	type question struct {
		id string
	}

	db, err := provider.ProvideSQL()
	if err != nil {
		log.Fatalf("failed provider sql: %v", err)
	}

	ImportSQL(db, "../../testdata/question/integration_test_questions.sql")

	defer cleaner()

	scenarios := []scenario{
		{
			name: "success create one question without space id",
			payload: map[string]interface{}{
				"question": "yoo, this work ?",
				"spaceId":  nil,
			},
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode, "must success create one question without spaceId")

				var q question

				err := db.
					QueryRowContext(ctx, "SELECT id FROM questions WHERE author_id = $1", "f028ac5a-e4c9-442f-bf9a-86c024a79baa").
					Scan(&q.id)

				assert.Nil(t, err)
			},
		},
		{
			name: "success create one question with space id",
			payload: map[string]interface{}{
				"question": "yoo, this work ?",
				"spaceId":  "a53152d7-2d24-42e1-a55f-649e87349ffa",
			},
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode, "must success create one question without spaceId")

				type question struct {
					id string
				}

				var q question

				err := db.
					QueryRowContext(ctx,
						"SELECT id FROM questions WHERE author_id = $1 AND space_id = $2",
						"f028ac5a-e4c9-442f-bf9a-86c024a79baa", "a53152d7-2d24-42e1-a55f-649e87349ffa",
					).
					Scan(&q.id)

				assert.Nil(t, err)
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

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/%s", url, "questions/"), bytes.NewReader(payload))
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
