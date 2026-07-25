package transaction_ledger

import (
	"testing"

	"github.com/cchristian77/wallet-service/domain"
	"github.com/cchristian77/wallet-service/domain/enums"
	m "github.com/cchristian77/wallet-service/mock"
	"github.com/cchristian77/wallet-service/request"
	"github.com/cchristian77/wallet-service/response"
	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"github.com/cchristian77/wallet-service/util"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (suite *TransactionLedgerServiceTestSuite) Test_Transfer() {
	fromWallet := m.InitWalletDomain(1, 10000)
	toWallet := m.InitWalletDomain(2, 5000)
	transaction := m.InitTransactionDomain()

	var expected *response.Transfer

	var input *request.Transfer

	testCases := []struct {
		name          string
		prepareMock   func()
		wantErr       bool
		expectedError error
	}{
		{
			name: "success",
			prepareMock: func() {
				suite.sqlMock.ExpectBegin()

				input = &request.Transfer{
					From:           1,
					To:             2,
					Amount:         1000,
					IdempotencyKey: "TRX-TEST-1",
				}

				expected = &response.Transfer{
					TransactionID: input.IdempotencyKey,
					FromBalance:   9000,
					ToBalance:     6000,
					Amount:        input.Amount,
					Status:        enums.TransactionStatusSuccess.String(),
				}

				suite.repo.EXPECT().FindTransactionByTransactionID(gomock.Any(), gomock.Eq(input.IdempotencyKey)).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).
					Return(transaction, nil).
					Times(1)
				suite.repo.EXPECT().FindWalletByIDForUpdate(gomock.Any(), gomock.Eq(uint64(1))).
					Return(m.InitWalletDomain(1, 10000), nil).
					Times(1)
				suite.repo.EXPECT().FindWalletByIDForUpdate(gomock.Any(), gomock.Eq(uint64(2))).
					Return(m.InitWalletDomain(2, 5000), nil).
					Times(1)
				suite.repo.EXPECT().CreateTransactionLedger(gomock.Any(), gomock.Any()).
					Return(&domain.TransactionLedger{}, nil).
					Times(2)
				suite.repo.EXPECT().UpdateWalletBalance(gomock.Any(), gomock.Eq(uint64(1)), gomock.Eq(int64(9000))).
					Return(nil).
					Times(1)
				suite.repo.EXPECT().UpdateWalletBalance(gomock.Any(), gomock.Eq(uint64(2)), gomock.Eq(int64(6000))).
					Return(nil).
					Times(1)
				suite.repo.EXPECT().UpdateTransactionStatus(gomock.Any(), gomock.Eq(transaction.ID), gomock.Eq(enums.TransactionStatusSuccess.String())).
					Return(nil).
					Times(1)

				suite.sqlMock.ExpectCommit()
			},
		},
		{
			name: "same wallet",
			prepareMock: func() {
				input = &request.Transfer{
					From:           1,
					To:             1,
					Amount:         1000,
					IdempotencyKey: "TRX-TEST-1",
				}

				suite.repo.EXPECT().FindTransactionByTransactionID(gomock.Any(), gomock.Eq(input.IdempotencyKey)).
					Times(0)
				suite.repo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).
					Times(0)
				suite.repo.EXPECT().FindWalletByIDForUpdate(gomock.Any(), gomock.Eq(uint64(1))).
					Times(0)
			},
			wantErr:       true,
			expectedError: sameWalletErr,
		},
		{
			name: "transaction already processed",
			prepareMock: func() {
				input = &request.Transfer{
					From:           1,
					To:             2,
					Amount:         1000,
					IdempotencyKey: "TRX-TEST-1",
				}

				suite.sqlMock.ExpectBegin()

				suite.repo.EXPECT().FindTransactionByTransactionID(gomock.Any(), gomock.Eq(input.IdempotencyKey)).
					Return(transaction, nil).
					Times(1)
				suite.repo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().FindWalletByIDForUpdate(gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().CreateTransactionLedger(gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().UpdateWalletBalance(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().UpdateTransactionStatus(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				suite.sqlMock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: transactionAlreadyProcessedErr,
		},
		{
			name: "from wallet not found",
			prepareMock: func() {
				input = &request.Transfer{
					From:           1,
					To:             2,
					Amount:         1000,
					IdempotencyKey: "TRX-TEST-1",
				}

				suite.sqlMock.ExpectBegin()

				suite.repo.EXPECT().FindTransactionByTransactionID(gomock.Any(), gomock.Eq(input.IdempotencyKey)).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).
					Return(transaction, nil).
					Times(1)
				suite.repo.EXPECT().FindWalletByIDForUpdate(gomock.Any(), gomock.Eq(uint64(1))).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().FindWalletByIDForUpdate(gomock.Any(), gomock.Eq(uint64(2))).Times(0)
				suite.repo.EXPECT().CreateTransactionLedger(gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().UpdateWalletBalance(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().UpdateTransactionStatus(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				suite.sqlMock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: sharedErrs.New(sharedErrs.ErrKindDataNotFound, "Wallet with id %d not found", uint64(1)),
		},
		{
			name: "to wallet not found",
			prepareMock: func() {
				input = &request.Transfer{
					From:           1,
					To:             2,
					Amount:         1000,
					IdempotencyKey: "TRX-TEST-1",
				}

				suite.sqlMock.ExpectBegin()

				suite.repo.EXPECT().FindTransactionByTransactionID(gomock.Any(), gomock.Eq(input.IdempotencyKey)).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).
					Return(transaction, nil).
					Times(1)
				suite.repo.EXPECT().FindWalletByIDForUpdate(gomock.Any(), gomock.Eq(uint64(1))).
					Return(fromWallet, nil).
					Times(1)
				suite.repo.EXPECT().FindWalletByIDForUpdate(gomock.Any(), gomock.Eq(uint64(2))).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().CreateTransactionLedger(gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().UpdateWalletBalance(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().UpdateTransactionStatus(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				suite.sqlMock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: sharedErrs.New(sharedErrs.ErrKindDataNotFound, "Wallet with id %d not found", uint64(1)),
		},
		{
			name: "insufficient balance",
			prepareMock: func() {
				input = &request.Transfer{
					From:           1,
					To:             2,
					Amount:         1000,
					IdempotencyKey: "TRX-TEST-1",
				}

				suite.sqlMock.ExpectBegin()

				suite.repo.EXPECT().FindTransactionByTransactionID(gomock.Any(), gomock.Eq(input.IdempotencyKey)).
					Return(nil, sharedErrs.NotFoundErr).
					Times(1)
				suite.repo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).
					Return(transaction, nil).
					Times(1)
				suite.repo.EXPECT().FindWalletByIDForUpdate(gomock.Any(), gomock.Eq(uint64(1))).
					Return(m.InitWalletDomain(1, 100), nil).
					Times(1)
				suite.repo.EXPECT().FindWalletByIDForUpdate(gomock.Any(), gomock.Eq(uint64(2))).
					Return(toWallet, nil).
					Times(1)
				suite.repo.EXPECT().CreateTransactionLedger(gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().UpdateWalletBalance(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				suite.repo.EXPECT().UpdateTransactionStatus(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				suite.sqlMock.ExpectRollback()
			},
			wantErr:       true,
			expectedError: insufficientBalanceErr,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Arrange
			suite.Before(t)
			defer suite.After(t)
			tc.prepareMock()

			// Act
			result, err := suite.transactionLedgerService.Transfer(suite.ctx, input)

			// Assert
			assert.Equal(t, tc.wantErr, err != nil, "error expected %v, but actual: %v", tc.wantErr, err)
			if tc.wantErr {
				assert.Empty(t, result)
				assert.Error(t, err)
			} else {
				assert.NotEmpty(t, result)
				if err = util.CompareData(result, expected, 1); err != nil {
					t.Errorf("error on comparing data : %v", err)
				}
			}

			assert.NoError(t, suite.sqlMock.ExpectationsWereMet())
		})
	}
}
