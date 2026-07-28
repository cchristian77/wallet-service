package memstore

import (
	"context"
	"sync"
	"time"

	"github.com/cchristian77/wallet-service/domain"
	"github.com/cchristian77/wallet-service/domain/enums"
	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
)

// MemStore is an in-memory, concurrency-safe Repository. It is the zero-dependency
// backend used for tests and for running without Postgres.
type MemStore struct {
	mu sync.RWMutex

	users map[uint64]*domain.User

	wallets      map[uint64]*domain.Wallet
	walletsByUID map[uint64]uint64 // userID -> walletID

	transactions map[uint64]*domain.Transaction
	txByKey      map[string]uint64 // idempotency key -> transaction ID

	ledgers      map[uint64]*domain.TransactionLedger
	ledgersByTxn map[uint64][]uint64 // transaction ID -> ledger IDs

	nextUserID   uint64
	nextWalletID uint64
	nextTxnID    uint64
	nextLedgerID uint64
}

// NewMemStore returns an empty in-memory store.
func NewMemStore() *MemStore {
	return &MemStore{
		users:        make(map[uint64]*domain.User),
		wallets:      make(map[uint64]*domain.Wallet),
		walletsByUID: make(map[uint64]uint64),
		transactions: make(map[uint64]*domain.Transaction),
		txByKey:      make(map[string]uint64),
		ledgers:      make(map[uint64]*domain.TransactionLedger),
		ledgersByTxn: make(map[uint64][]uint64),
		nextUserID:   1,
		nextWalletID: 1,
		nextTxnID:    1,
		nextLedgerID: 1,
	}
}

// NewSeededMemStore returns a MemStore preloaded from JSON seed files
// (users, wallets, transactions, transaction_ledgers).
func NewSeededMemStore() (*MemStore, error) {
	s := NewMemStore()
	if err := s.loadSeed(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *MemStore) FindWalletByID(ctx context.Context, id uint64) (*domain.Wallet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w, ok := s.wallets[id]
	if !ok {
		return nil, sharedErrs.NotFoundErr
	}
	return cloneWallet(w), nil
}

func (s *MemStore) FindWalletByUserID(ctx context.Context, userID uint64) (*domain.Wallet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	walletID, ok := s.walletsByUID[userID]
	if !ok {
		return nil, sharedErrs.NotFoundErr
	}
	w, ok := s.wallets[walletID]
	if !ok {
		return nil, sharedErrs.NotFoundErr
	}
	return cloneWallet(w), nil
}

func (s *MemStore) FindWalletByIDForUpdate(ctx context.Context, id uint64) (*domain.Wallet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.wallets[id]
	if !ok {
		return nil, sharedErrs.NotFoundErr
	}
	return cloneWallet(w), nil
}

func (s *MemStore) UpdateWalletBalance(ctx context.Context, walletID uint64, balance int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.wallets[walletID]
	if !ok {
		return sharedErrs.NotFoundErr
	}
	w.Balance = balance
	w.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *MemStore) FindTransactionByTransactionID(ctx context.Context, transactionID string) (*domain.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.txByKey[transactionID]
	if !ok {
		return nil, sharedErrs.NotFoundErr
	}
	t, ok := s.transactions[id]
	if !ok {
		return nil, sharedErrs.NotFoundErr
	}
	return cloneTransaction(t), nil
}

func (s *MemStore) CreateTransaction(ctx context.Context, data *domain.Transaction) (*domain.Transaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.txByKey[data.TransactionID]; exists {
		return nil, sharedErrs.ConflictErr
	}

	now := time.Now().UTC()
	created := cloneTransaction(data)
	created.ID = s.nextTxnID
	s.nextTxnID++
	if created.CreatedAt.IsZero() {
		created.CreatedAt = now
	}
	created.UpdatedAt = now

	s.transactions[created.ID] = created
	s.txByKey[created.TransactionID] = created.ID
	return cloneTransaction(created), nil
}

func (s *MemStore) UpdateTransactionStatus(ctx context.Context, id uint64, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.transactions[id]
	if !ok {
		return sharedErrs.NotFoundErr
	}
	t.Status = enums.TransactionStatus(status)
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *MemStore) FindTransactionLedgersByTransactionID(ctx context.Context, transactionID uint64) ([]*domain.TransactionLedger, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := s.ledgersByTxn[transactionID]
	result := make([]*domain.TransactionLedger, 0, len(ids))
	for _, id := range ids {
		if l, ok := s.ledgers[id]; ok {
			result = append(result, cloneLedger(l))
		}
	}
	return result, nil
}

func (s *MemStore) CreateTransactionLedger(ctx context.Context, data *domain.TransactionLedger) (*domain.TransactionLedger, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.transactions[data.TransactionID]; !ok {
		return nil, sharedErrs.NotFoundErr
	}
	if _, ok := s.wallets[data.WalletID]; !ok {
		return nil, sharedErrs.NotFoundErr
	}

	now := time.Now().UTC()
	created := cloneLedger(data)
	created.ID = s.nextLedgerID
	s.nextLedgerID++
	if created.CreatedAt.IsZero() {
		created.CreatedAt = now
	}
	created.UpdatedAt = now

	s.ledgers[created.ID] = created
	s.ledgersByTxn[created.TransactionID] = append(s.ledgersByTxn[created.TransactionID], created.ID)
	return cloneLedger(created), nil
}

func cloneWallet(w *domain.Wallet) *domain.Wallet {
	if w == nil {
		return nil
	}
	cp := *w
	return &cp
}

func cloneTransaction(t *domain.Transaction) *domain.Transaction {
	if t == nil {
		return nil
	}
	cp := *t
	return &cp
}

func cloneLedger(l *domain.TransactionLedger) *domain.TransactionLedger {
	if l == nil {
		return nil
	}
	cp := *l
	return &cp
}
