package cassandra

import (
	"errors"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/stellar/federation/config"
	"github.com/stellar/federation/db"
)

type CassandraDriver struct {
	session *gocql.Session
}

func (d *CassandraDriver) Init(config config.ConfigDatabase) (err error) {
	cluster := gocql.NewCluster(config.Cluster...)
	cluster.Keyspace = config.Keyspace
	// Optional
	if config.ProtoVersion != nil {
		cluster.ProtoVersion = *config.ProtoVersion
	}
	d.session, err = cluster.CreateSession()
	return
}

func (d *CassandraDriver) GetByStellarAddress(name, query string) (*db.FederationRecord, error) {
	var record db.FederationRecord
	row := make(map[string]interface{})

	err := d.session.Query(query, name).MapScan(row)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}

	record.AccountId = row["id"].(string)
	memo, exist := row["memo"]
	if exist {
		switch memo := memo.(type) {
		case string:
			record.Memo = memo
		case int:
			record.Memo = strconv.Itoa(memo)
		default:
			return nil, errors.New("Unknown memo type")
		}

		// TODO
		record.MemoType = "text"
	}

	return &record, nil
}

func (d *CassandraDriver) GetByAccountId(accountId, query string) (*db.ReverseFederationRecord, error) {
	var record db.ReverseFederationRecord
	err := d.session.Query(query, accountId).Scan(&record.Name)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}
