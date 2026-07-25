package repository

import (
	"context"
	"testing"

	m "github.com/cchristian77/wallet-service/mock"
	"github.com/cchristian77/wallet-service/util/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestNewRepository_Success(t *testing.T) {
	writeDB, _, err := m.NewMockDB()
	if err != nil {
		t.Fatal(err)
	}

	repository := NewRepository(writeDB)

	assert.NotNil(t, repository)
	assert.Implements(t, (*Repository)(nil), repository)
}

type RepositoryTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	repo    Repository
	ctx     context.Context
}

func (suite *RepositoryTestSuite) Before(t *testing.T) {
	writerDB, sqlMock, err := m.NewMockDB()
	if err != nil {
		t.Fatal(err)
	}

	logger.Initialize()

	suite.sqlMock = sqlMock
	suite.repo = NewRepository(writerDB)
}

func TestRepositoryTestSuite(t *testing.T) {
	testSuite := new(RepositoryTestSuite)
	testSuite.ctx = context.Background()

	suite.Run(t, testSuite)
}
