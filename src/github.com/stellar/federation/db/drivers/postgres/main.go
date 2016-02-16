package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stellar/federation/db"
)

type PostgresDriver struct {
	database *sqlx.DB
}

func (d *PostgresDriver) Init(url string) (err error) {
	d.database, err = sqlx.Connect("postgres", url)
	return
}

func (d *PostgresDriver) GetByStellarAddress(name, query string) (*db.FederationRecord, error) {
	var record db.FederationRecord
	err := d.database.Get(&record, query, name)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (d *PostgresDriver) GetByAccountId(accountId, query string) (*db.ReverseFederationRecord, error) {
	var record db.ReverseFederationRecord
	err := d.database.Get(&record, query, accountId)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
