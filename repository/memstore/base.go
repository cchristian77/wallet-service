package memstore

import (
	"context"
	"sync"

	"github.com/cchristian77/wallet-service/domain"
)

// MemStore is an in-memory, concurrency-safe Repository.
// Begin holds the store mutex for the unit of work; Commit/Rollback releases it.
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

// NewSeededMemStore returns a MemStore preloaded from JSON seed files.
func NewSeededMemStore() (*MemStore, error) {
	s := NewMemStore()

	if err := s.loadSeed(); err != nil {
		return nil, err
	}

	return s, nil
}

// rLock acquires a shared lock unless a Begin()-started tx already holds the store.
func (s *MemStore) rLock(ctx context.Context) {
	if inTx(ctx) {
		return
	}
	s.mu.RLock()
}

func (s *MemStore) rUnlock(ctx context.Context) {
	if inTx(ctx) {
		return
	}
	s.mu.RUnlock()
}

// wLock acquires an exclusive lock unless a Begin()-started tx already holds the store.
func (s *MemStore) wLock(ctx context.Context) {
	if inTx(ctx) {
		return
	}
	s.mu.Lock()
}

func (s *MemStore) wUnlock(ctx context.Context) {
	if inTx(ctx) {
		return
	}
	s.mu.Unlock()
}
