package memstore

import (
	"context"
	"database/sql"

	"github.com/cchristian77/wallet-service/domain"
	"github.com/cchristian77/wallet-service/shared/external/database"
)

type txKey struct{}

type Transactor struct {
	store    *MemStore
	snapshot *storeSnapshot
	done     bool
}

// Begin starts a unit of work and holds the store mutex until Commit or Rollback.
func (s *MemStore) Begin(ctx context.Context) (context.Context, database.Tx) {
	s.mu.Lock()

	tx := &Transactor{
		store:    s,
		snapshot: s.takeSnapshot(),
	}

	return context.WithValue(ctx, txKey{}, tx), tx
}

func (t *Transactor) Commit() error {
	if t.done {
		return sql.ErrTxDone
	}

	t.done = true
	t.store.mu.Unlock()

	return nil
}

func (t *Transactor) Rollback() error {
	if t.done {
		return sql.ErrTxDone
	}

	t.done = true
	t.store.restoreSnapshot(t.snapshot)
	t.store.mu.Unlock()

	return nil
}

func inTx(ctx context.Context) bool {
	tx, _ := ctx.Value(txKey{}).(*Transactor)
	return tx != nil && !tx.done
}

type storeSnapshot struct {
	users        map[uint64]*domain.User
	wallets      map[uint64]*domain.Wallet
	walletsByUID map[uint64]uint64
	transactions map[uint64]*domain.Transaction
	txByKey      map[string]uint64
	ledgers      map[uint64]*domain.TransactionLedger
	ledgersByTxn map[uint64][]uint64
	nextUserID   uint64
	nextWalletID uint64
	nextTxnID    uint64
	nextLedgerID uint64
}

// takeSnapshot deep-copies store state. Caller must hold s.mu.
func (s *MemStore) takeSnapshot() *storeSnapshot {
	snap := &storeSnapshot{
		users:        make(map[uint64]*domain.User, len(s.users)),
		wallets:      make(map[uint64]*domain.Wallet, len(s.wallets)),
		walletsByUID: make(map[uint64]uint64, len(s.walletsByUID)),
		transactions: make(map[uint64]*domain.Transaction, len(s.transactions)),
		txByKey:      make(map[string]uint64, len(s.txByKey)),
		ledgers:      make(map[uint64]*domain.TransactionLedger, len(s.ledgers)),
		ledgersByTxn: make(map[uint64][]uint64, len(s.ledgersByTxn)),
		nextUserID:   s.nextUserID,
		nextWalletID: s.nextWalletID,
		nextTxnID:    s.nextTxnID,
		nextLedgerID: s.nextLedgerID,
	}

	for k, v := range s.users {
		snap.users[k] = cloneUser(v)
	}
	for k, v := range s.wallets {
		snap.wallets[k] = cloneWallet(v)
	}
	for k, v := range s.walletsByUID {
		snap.walletsByUID[k] = v
	}
	for k, v := range s.transactions {
		snap.transactions[k] = cloneTransaction(v)
	}
	for k, v := range s.txByKey {
		snap.txByKey[k] = v
	}
	for k, v := range s.ledgers {
		snap.ledgers[k] = cloneLedger(v)
	}
	for k, ids := range s.ledgersByTxn {
		snap.ledgersByTxn[k] = append([]uint64(nil), ids...)
	}

	return snap
}

// restoreSnapshot replaces store state with snap. Caller must hold s.mu.
func (s *MemStore) restoreSnapshot(snap *storeSnapshot) {
	s.users = snap.users
	s.wallets = snap.wallets
	s.walletsByUID = snap.walletsByUID
	s.transactions = snap.transactions
	s.txByKey = snap.txByKey
	s.ledgers = snap.ledgers
	s.ledgersByTxn = snap.ledgersByTxn
	s.nextUserID = snap.nextUserID
	s.nextWalletID = snap.nextWalletID
	s.nextTxnID = snap.nextTxnID
	s.nextLedgerID = snap.nextLedgerID
}
