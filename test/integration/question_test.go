package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/rizface/quora/account/value"
	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationTestSuite) TestCreateNewQuestion() {
	type question struct {
		id string
	}

	authenticated, err := value.NewAuthenticated(value.AccountEntity{
		Id:       "f028ac5a-e4c9-442f-bf9a-86c024a79baa",
		Username: "testlogin",
		Email:    "testlogin@gmail.com",
	})
	if err != nil {
		suite.Error(err)
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
			url, err := suite.services.quora.Endpoint(suite.ctx, "")
			if err != nil {
				t.Error(err)
			}

			r := requester{
				url:     fmt.Sprintf("http://%s/%s", url, "questions/"),
				payload: s.payload,
				method:  http.MethodPost,
				headers: map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", authenticated.Tokens[0].Value),
				},
			}

			resp, err := r.do()
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			if s.checkExpectation != nil {
				s.checkExpectation(t, resp)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestVoteAnswer() {
	type scenario struct {
		name             string
		answerId         string
		voteType         string
		checPreTest      func(t *testing.T)
		checkExpectation func(t *testing.T, resp *http.Response)
	}

	type answer struct {
		upvote   int
		downvote int
	}

	authenticated, err := value.NewAuthenticated(value.AccountEntity{
		Id:       "f028ac5a-e4c9-442f-bf9a-86c024a79baa",
		Username: "testlogin",
		Email:    "testlogin@gmail.com",
	})
	if err != nil {
		suite.Error(err)
	}

	ImportSQL(suite.db, "../../testdata/question/integration_test_questions.sql")

	scenarios := []scenario{
		{
			name:     "success upvote one answer",
			answerId: "4b9ef364-0d6a-4f60-a169-39b1d076c65d",
			voteType: "upvote",
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var (
					question = answer{}
					query    = `
						select upvote from answers where id = $1
					`
				)

				err := suite.db.QueryRowContext(suite.ctx, query, "4b9ef364-0d6a-4f60-a169-39b1d076c65d").Scan(&question.upvote)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, 1, question.upvote)
			},
		},
		{
			name:     "success downvote one question",
			answerId: "4b9ef364-0d6a-4f60-a169-39b1d076c65f",
			voteType: "downvote",
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var (
					question = answer{}
					query    = `
						select downvote from answers where id = $1
					`
				)

				err := suite.db.QueryRowContext(suite.ctx, query, "4b9ef364-0d6a-4f60-a169-39b1d076c65f").Scan(&question.downvote)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, 1, question.downvote)
			},
		},
		{
			name:     "invalid vote type",
			answerId: "4b9ef364-0d6a-4f60-a169-39b1d076c65d",
			voteType: "invalidvote",
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:     "answer not found",
			answerId: "a53152d7-2d24-42e1-a55f-649e87349ffa",
			voteType: "upvote",
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name:     "spam upvote or downvote must be ignored",
			answerId: "4b9ef364-0d6a-4f60-a169-39b1d076c65e",
			voteType: "upvote",
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var (
					question = answer{}
					query    = `
						select upvote from answers where id = $1
					`
				)

				err := suite.db.QueryRowContext(suite.ctx, query, "4b9ef364-0d6a-4f60-a169-39b1d076c65e").Scan(&question.upvote)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, 1, question.upvote)
			},
		},
		{
			name:     "existing vote must be deleted if client send the opposite vote",
			answerId: "4b9ef364-0d6a-4f60-a169-39b1d076c65b",
			voteType: "upvote",
			checPreTest: func(t *testing.T) {
				var (
					question = answer{}
					query    = `
						select upvote, downvote from answers where id = $1
					`
				)

				err := suite.db.QueryRowContext(suite.ctx, query, "4b9ef364-0d6a-4f60-a169-39b1d076c65b").Scan(
					&question.upvote,
					&question.downvote,
				)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, 0, question.upvote, "upvote must be 0 before re-vote the question")
				assert.Equal(t, 1, question.downvote, "downvote must be 1 before re-vote the question")
			},
			checkExpectation: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var (
					question = answer{}
					query    = `
						select upvote,downvote from answers where id = $1
					`
				)

				err := suite.db.QueryRowContext(suite.ctx, query, "4b9ef364-0d6a-4f60-a169-39b1d076c65b").Scan(
					&question.upvote,
					&question.downvote,
				)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, 1, question.upvote, "upvote must be 1 after re-vote the question")
				assert.Equal(t, 0, question.downvote, "downvote must be 0 after re-vote the question")
			},
		},
	}

	t := suite.T()

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			if s.checPreTest != nil {
				s.checPreTest(t)
			}

			url, err := suite.services.quora.Endpoint(suite.ctx, "")
			if err != nil {
				t.Error(err)
			}

			r := requester{
				url: fmt.Sprintf("http://%s/%s/%s/vote", url, "answers", s.answerId),
				payload: map[string]interface{}{
					"type": s.voteType,
				},
				method: http.MethodPatch,
				headers: map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", authenticated.Tokens[0].Value),
				},
			}

			resp, err := r.do()
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()

			if s.checkExpectation != nil {
				s.checkExpectation(t, resp)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestCreateAnswer() {
	type (
		scenario struct {
			name             string
			payload          map[string]interface{}
			checkExpectation func(resp *http.Response)
		}

		answer struct {
			id string
		}
	)

	authenticated, err := value.NewAuthenticated(value.AccountEntity{
		Id:       "f028ac5a-e4c9-442f-bf9a-86c024a79baa",
		Username: "testlogin",
		Email:    "testlogin@gmail.com",
	})
	if err != nil {
		suite.Error(err)
	}

	ImportSQL(suite.db, "../../testdata/question/integration_test_questions.sql")

	scenarios := []scenario{
		{
			name: "success create answer for question",
			payload: map[string]interface{}{
				"answer":     "this is good questions",
				"questionId": "4b9ef364-0d6a-4f60-a169-39b1d076c63a",
			},
			checkExpectation: func(resp *http.Response) {
				suite.Equal(http.StatusOK, resp.StatusCode)

				var (
					answer = answer{}
					query  = `
						select id from answers where question_id = $1
					`
				)

				err := suite.db.QueryRowContext(suite.ctx, query, "4b9ef364-0d6a-4f60-a169-39b1d076c63a").Scan(&answer.id)
				suite.NoError(err)
			},
		},
		{
			name: "failed create answer for not found question",
			payload: map[string]interface{}{
				"answer":     "this is good questions",
				"questionId": "4b9ef364-0d6a-4f60-a169-39b1d076c62a",
			},
			checkExpectation: func(resp *http.Response) {
				suite.Equal(http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name: "failed create empty answer",
			payload: map[string]interface{}{
				"answer":     "",
				"questionId": "4b9ef364-0d6a-4f60-a169-39b1d076c63a",
			},
			checkExpectation: func(resp *http.Response) {
				suite.Equal(http.StatusBadRequest, resp.StatusCode)

				var (
					answer = answer{}
					query  = `
						select id from answers where question_id = $1
					`
				)

				err := suite.db.QueryRowContext(suite.ctx, query, "4b9ef364-0d6a-4f60-a169-39b1d076c63a").Scan(&answer.id)
				suite.NoError(err)
			},
		},
	}

	for _, s := range scenarios {
		suite.Run(s.name, func() {
			url, err := suite.services.quora.Endpoint(suite.ctx, "")
			if err != nil {
				suite.Error(err)
			}

			r := requester{
				url:     fmt.Sprintf("http://%s/%s", url, "answers"),
				payload: s.payload,
				method:  http.MethodPost,
				headers: map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", authenticated.Tokens[0].Value),
				},
			}

			resp, err := r.do()
			if err != nil {
				suite.Error(err)
			}
			defer resp.Body.Close()

			if s.checkExpectation != nil {
				s.checkExpectation(resp)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestDeleteQuestion() {
	type (
		scenario struct {
			name             string
			questionId       string
			token            string
			checkExpectation func(resp *http.Response)
		}

		question struct {
			id string
		}
	)

	var (
		// mimick user in database
		users = map[string]value.AccountEntity{
			"user1": {
				Id:       "f028ac5a-e4c9-442f-bf9a-86c024a79baa",
				Username: "testlogin",
				Email:    "testlogin@gmail.com",
			},
			"user2": {
				Id:       "f028ac5a-e4c9-442f-bf9a-86c024a79bac",
				Username: "testdelete",
				Email:    "testdelete@gmail.com",
			},
		}

		usersToken = map[string]string{}
	)

	for k, v := range users {
		authenticated, err := value.NewAuthenticated(v)
		if err != nil {
			suite.T().Fatal(err)
		}

		usersToken[k] = authenticated.Tokens[0].Value
	}

	ImportSQL(suite.db, "../../testdata/question/integration_test_questions.sql")

	scenarios := []scenario{
		{
			name:       "success delete one question",
			questionId: "4b9ef364-0d6a-4f60-a169-39b1d076c63c",
			token:      usersToken["user1"],
			checkExpectation: func(resp *http.Response) {
				suite.Equal(http.StatusOK, resp.StatusCode)

				question := question{}

				err := suite.db.QueryRow(`SELECT id FROM questions where id = $1`, "4b9ef364-0d6a-4f60-a169-39b1d076c63c").Scan(&question.id)

				suite.ErrorIs(err, sql.ErrNoRows)
			},
		},
		{
			name:       "failed delete one question (not found)",
			questionId: "4b9ef364-0d6a-4f60-a169-39b1d076c62b",
			token:      usersToken["user1"],
			checkExpectation: func(resp *http.Response) {
				suite.Equal(http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name:       "failed delete one question (not the author)",
			questionId: "4b9ef364-0d6a-4f60-a169-39b1d076c65e",
			token:      usersToken["user2"],
			checkExpectation: func(resp *http.Response) {
				suite.Equal(http.StatusUnauthorized, resp.StatusCode)
			},
		},
	}

	for _, s := range scenarios {
		suite.Run(s.name, func() {
			url, err := suite.services.quora.Endpoint(suite.ctx, "")
			if err != nil {
				suite.Error(err)
			}

			r := requester{
				url:     fmt.Sprintf("http://%s/%s/%s", url, "questions", s.questionId),
				payload: nil,
				method:  http.MethodDelete,
				headers: map[string]string{
					"Authorization": "Bearer " + s.token,
				},
			}

			resp, err := r.do()
			if err != nil {
				suite.T().Error(err)
			}

			defer resp.Body.Close()

			if s.checkExpectation != nil {
				s.checkExpectation(resp)
			}
		})
	}
}

func prinResponse(t *testing.T, w *http.Response) {
	body := map[string]interface{}{}

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	fmt.Println(body)
}
