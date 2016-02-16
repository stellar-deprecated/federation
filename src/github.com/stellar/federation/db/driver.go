package db

type Driver interface {
	Init(url string) (err error)
	GetByStellarAddress(name, query string) (*FederationRecord, error)
	GetByAccountId(accountId, query string) (*ReverseFederationRecord, error)
}
