package integration

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/rizface/quora/provider"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	ctx context.Context
	suite.Suite
	db       *sql.DB
	services services
	cleaner  func()
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.services, suite.cleaner = spawnServices(suite.ctx)

	db, err := provider.ProvideSQL()
	if err != nil {
		log.Fatalf("failed when start test suite: %v", err)
	}

	suite.db = db
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.cleaner()
}

func TestIntegrationTest(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
