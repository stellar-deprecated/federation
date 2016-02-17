package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stellar/federation/config"
	"github.com/stellar/federation/db"
)

type MysqlDriver struct {
	database *sqlx.DB
}

func (d *MysqlDriver) Init(config config.ConfigDatabase) (err error) {
	d.database, err = sqlx.Connect("mysql", config.Url)
	return
}

func (d *MysqlDriver) GetByStellarAddress(name, query string) (*db.FederationRecord, error) {
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

func (d *MysqlDriver) GetByAccountId(accountId, query string) (*db.ReverseFederationRecord, error) {
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
