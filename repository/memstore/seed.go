package memstore

import (
	"embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cchristian77/wallet-service/domain"
	"github.com/cchristian77/wallet-service/domain/enums"
)

//go:embed users.json wallets.json transactions.json transaction_ledgers.json
var seedFS embed.FS

type userSeed struct {
	ID       uint64 `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type walletSeed struct {
	ID      uint64 `json:"id"`
	UserID  uint64 `json:"user_id"`
	Balance int64  `json:"balance"`
}

type transactionSeed struct {
	ID            uint64 `json:"id"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

type transactionLedgerSeed struct {
	ID            uint64 `json:"id"`
	TransactionID uint64 `json:"transaction_id"`
	WalletID      uint64 `json:"wallet_id"`
	Ledger        string `json:"ledger"`
	Reference     string `json:"reference"`
	Amount        int64  `json:"amount"`
}

func (s *MemStore) loadSeed() error {
	if err := s.loadUsers(); err != nil {
		return err
	}
	if err := s.loadWallets(); err != nil {
		return err
	}
	if err := s.loadTransactions(); err != nil {
		return err
	}
	if err := s.loadTransactionLedgers(); err != nil {
		return err
	}

	return nil
}

func (s *MemStore) loadUsers() error {
	var rows []userSeed
	if err := decodeSeed("users.json", &rows); err != nil {
		return err
	}

	now := time.Now()
	var lastID uint64

	for _, row := range rows {
		s.users[row.ID] = &domain.User{
			BaseModel: domain.BaseModel{
				ID:        row.ID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			FullName: row.FullName,
			Email:    row.Email,
			Password: row.Password,
		}
		if row.ID > lastID {
			lastID = row.ID
		}
	}
	s.nextUserID = lastID + 1

	return nil
}

func (s *MemStore) loadWallets() error {
	var rows []walletSeed
	if err := decodeSeed("wallets.json", &rows); err != nil {
		return err
	}

	now := time.Now()
	var lastID uint64

	for _, row := range rows {
		s.wallets[row.ID] = &domain.Wallet{
			BaseModel: domain.BaseModel{
				ID:        row.ID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			UserID:  row.UserID,
			Balance: row.Balance,
		}
		s.walletsByUID[row.UserID] = row.ID
		if row.ID > lastID {
			lastID = row.ID
		}
	}
	s.nextWalletID = lastID + 1

	return nil
}

func (s *MemStore) loadTransactions() error {
	var rows []transactionSeed
	if err := decodeSeed("transactions.json", &rows); err != nil {
		return err
	}

	now := time.Now()
	var lastID uint64

	for _, row := range rows {
		s.transactions[row.ID] = &domain.Transaction{
			BaseModel: domain.BaseModel{
				ID:        row.ID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			TransactionID: row.TransactionID,
			Status:        enums.TransactionStatus(row.Status),
		}
		s.txByKey[row.TransactionID] = row.ID
		if row.ID > lastID {
			lastID = row.ID
		}
	}
	s.nextTxnID = lastID + 1

	return nil
}

func (s *MemStore) loadTransactionLedgers() error {
	var rows []transactionLedgerSeed
	if err := decodeSeed("transaction_ledgers.json", &rows); err != nil {
		return err
	}

	now := time.Now()
	var lastID uint64

	for _, row := range rows {
		s.ledgers[row.ID] = &domain.TransactionLedger{
			BaseModel: domain.BaseModel{
				ID:        row.ID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			TransactionID: row.TransactionID,
			WalletID:      row.WalletID,
			Ledger:        enums.TransactionLedgerType(row.Ledger),
			Reference:     enums.TransactionLedgerReference(row.Reference),
			Amount:        row.Amount,
		}
		s.ledgersByTxn[row.TransactionID] = append(s.ledgersByTxn[row.TransactionID], row.ID)
		if row.ID > lastID {
			lastID = row.ID
		}
	}
	s.nextLedgerID = lastID + 1

	return nil
}

func decodeSeed(name string, dest any) error {
	b, err := seedFS.ReadFile(name)
	if err != nil {
		return fmt.Errorf("memstore: read %s: %w", name, err)
	}
	if err := json.Unmarshal(b, dest); err != nil {
		return fmt.Errorf("memstore: parse %s: %w", name, err)
	}

	return nil
}
