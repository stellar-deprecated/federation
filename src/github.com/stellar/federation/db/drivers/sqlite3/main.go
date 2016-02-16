package sqlite3

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stellar/federation/db"
)

type Sqlite3Driver struct {
	database *sqlx.DB
}

func (d *Sqlite3Driver) Init(url string) (err error) {
	d.database, err = sqlx.Connect("sqlite3", url)
	return
}

func (d *Sqlite3Driver) GetByStellarAddress(name, query string) (*db.FederationRecord, error) {
	var record db.FederationRecord
	err := d.database.Get(&record, query, name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}

func (d *Sqlite3Driver) GetByAccountId(accountId, query string) (*db.ReverseFederationRecord, error) {
	var record db.ReverseFederationRecord
	err := d.database.Get(&record, query, accountId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}
