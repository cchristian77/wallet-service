package transaction_ledger

import (
	"context"
	"testing"

	m "github.com/cchristian77/wallet-service/mock"
	"github.com/cchristian77/wallet-service/util/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestNewService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := m.NewMockRepository(ctrl)
	writeDB, _, err := m.NewMockDB()
	if err != nil {
		t.Fatal(err)
	}

	transactionLedgerService, err := NewService(repoMock, writeDB)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, transactionLedgerService)
	assert.Implements(t, (*Service)(nil), transactionLedgerService)
}

type TransactionLedgerServiceTestSuite struct {
	suite.Suite
	repo    *m.MockRepository
	writeDB *gorm.DB
	sqlMock sqlmock.Sqlmock
	ctx     context.Context

	transactionLedgerService Service
}

func (suite *TransactionLedgerServiceTestSuite) Before(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var err error

	logger.InitializeNop()

	suite.repo = m.NewMockRepository(ctrl)
	suite.writeDB, suite.sqlMock, err = m.NewMockDB()
	if err != nil {
		t.Fatal(err)
	}

	suite.transactionLedgerService, err = NewService(suite.repo, suite.writeDB)
	if err != nil {
		t.Fatal(err)
	}
}

func (suite *TransactionLedgerServiceTestSuite) After(t *testing.T) {}

func TestSuiteRunTransactionLedgerService(t *testing.T) {
	suite.Run(t, new(TransactionLedgerServiceTestSuite))
}
