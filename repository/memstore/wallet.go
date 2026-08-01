package memstore

import (
	"context"
	"time"

	"github.com/cchristian77/wallet-service/domain"
	sharedErrs "github.com/cchristian77/wallet-service/shared/errors"
	"github.com/cchristian77/wallet-service/util/logger"
)

func (s *MemStore) FindWalletByID(ctx context.Context, id uint64) (*domain.Wallet, error) {
	s.rLock(ctx)
	defer s.rUnlock(ctx)

	w, ok := s.wallets[id]
	if !ok {
		logger.Error(ctx, "[MEMSTORE] Failed on find wallet by id : not found")
		return nil, sharedErrs.NotFoundErr
	}

	return cloneWallet(w), nil
}

func (s *MemStore) FindWalletByUserID(ctx context.Context, userID uint64) (*domain.Wallet, error) {
	s.rLock(ctx)
	defer s.rUnlock(ctx)

	walletID, ok := s.walletsByUID[userID]
	if !ok {
		logger.Error(ctx, "[MEMSTORE] Failed on find wallet by user id : not found")
		return nil, sharedErrs.NotFoundErr
	}

	w, ok := s.wallets[walletID]
	if !ok {
		logger.Error(ctx, "[MEMSTORE] Failed on find wallet by user id : not found")
		return nil, sharedErrs.NotFoundErr
	}

	return cloneWallet(w), nil
}

func (s *MemStore) FindWalletByIDForUpdate(ctx context.Context, id uint64) (*domain.Wallet, error) {
	s.wLock(ctx)
	defer s.wUnlock(ctx)

	w, ok := s.wallets[id]
	if !ok {
		logger.Error(ctx, "[MEMSTORE] Failed on find wallet by id for update : not found")
		return nil, sharedErrs.NotFoundErr
	}

	return cloneWallet(w), nil
}

func (s *MemStore) UpdateWalletBalance(ctx context.Context, id uint64, balance int64) error {
	s.wLock(ctx)
	defer s.wUnlock(ctx)

	w, ok := s.wallets[id]
	if !ok {
		logger.Error(ctx, "[MEMSTORE] Failed on update wallet balance: not found")
		return sharedErrs.NotFoundErr
	}

	w.Balance = balance
	w.UpdatedAt = time.Now().UTC()

	return nil
}
