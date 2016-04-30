package db

import (
	"github.com/stellar/federation/config"
)

type Driver interface {
	Init(config config.ConfigDatabase) (err error)
	GetByStellarAddress(name, query string) (*FederationRecord, error)
	GetByAccountId(accountId, query string) (*ReverseFederationRecord, error)
}
