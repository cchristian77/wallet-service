package mock

import (
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cchristian77/wallet-service/domain"
	"github.com/cchristian77/wallet-service/domain/enums"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
This file provides functionality to create instances of the specified required structs for unit testing purposes.
This ensures that tests have consistent and predictable data without the need for creating these objects manually in each test case.
*/

/*
 * ============================= MOCKING =============================
 */

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("Error occurs when opening a stub database connection : %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	}), &gorm.Config{})

	if err != nil {
		return nil, nil, fmt.Errorf("Error occurs when opening gorm database : %v", err)
	}

	return gormDB, mock, err
}

/*
 * ============================= DOMAIN =============================
 */

func InitWalletDomain(id uint64, balance int64) *domain.Wallet {
	now := time.Now()

	return &domain.Wallet{
		BaseModel: domain.BaseModel{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:  id,
		Balance: balance,
	}
}

func InitTransactionDomain() *domain.Transaction {
	now := time.Now()

	return &domain.Transaction{
		BaseModel: domain.BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		TransactionID: "TRX-TEST-1",
		Status:        enums.TransactionStatusPending,
	}
}
