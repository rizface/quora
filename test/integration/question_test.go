package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationTestSuite) TestCreateNewQuestion() {
	type question struct {
		id string
	}

	ImportSQL(suite.db, "../../testdata/question/integration_test_questions.sql")

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

				err := suite.db.
					QueryRowContext(suite.ctx, "SELECT id FROM questions WHERE author_id = $1", "f028ac5a-e4c9-442f-bf9a-86c024a79baa").
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

				err := suite.db.
					QueryRowContext(suite.ctx,
						"SELECT id FROM questions WHERE author_id = $1 AND space_id = $2",
						"f028ac5a-e4c9-442f-bf9a-86c024a79baa", "a53152d7-2d24-42e1-a55f-649e87349ffa",
					).
					Scan(&q.id)

				assert.Nil(t, err)
			},
		},
	}

	t := suite.T()

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			client := &http.Client{}

			url, err := suite.services.quora.Endpoint(suite.ctx, "")
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

func (suite *IntegrationTestSuite) TestVoteQuestion() {
	type scenario struct {
		name             string
		questionId       string
		voteType         string
		checkExpectation func(t *testing.T, resp *http.Response)
	}

	type question struct {
		upvote   int
		downvote int
	}

	ImportSQL(suite.db, "../../testdata/question/integration_test_questions.sql")

	scenarios := []scenario{
		{
			name:       "success upvote one question",
			questionId: "4b9ef364-0d6a-4f60-a169-39b1d076c65d",
			voteType:   "upvote",
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var (
					question = question{}
					query    = `
						select upvote from questions limit 1
					`
				)

				err := suite.db.QueryRowContext(suite.ctx, query).Scan(&question.upvote)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, 1, question.upvote)
			},
		},
		{
			name:       "success downvote one question",
			questionId: "4b9ef364-0d6a-4f60-a169-39b1d076c65d",
			voteType:   "downvote",
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var (
					question = question{}
					query    = `
						select downvote from questions limit 1
					`
				)

				err := suite.db.QueryRowContext(suite.ctx, query).Scan(&question.downvote)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, 1, question.downvote)
			},
		},
		{
			name:       "invalid vote type",
			questionId: "4b9ef364-0d6a-4f60-a169-39b1d076c65d",
			voteType:   "invalidvote",
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:       "question not found",
			questionId: "a53152d7-2d24-42e1-a55f-649e87349ffa",
			voteType:   "upvote",
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			},
		},
	}

	t := suite.T()

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			client := &http.Client{}

			url, err := suite.services.quora.Endpoint(suite.ctx, "")
			if err != nil {
				t.Error(err)
			}

			payload, err := json.Marshal(map[string]interface{}{
				"type": s.voteType,
			})
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("http://%s/%s/%s/vote", url, "questions", s.questionId), bytes.NewReader(payload))
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
