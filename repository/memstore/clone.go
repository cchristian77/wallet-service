package memstore

import "github.com/cchristian77/wallet-service/domain"

func cloneUser(u *domain.User) *domain.User {
	if u == nil {
		return nil
	}
	cp := *u
	return &cp
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
