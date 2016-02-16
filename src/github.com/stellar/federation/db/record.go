package db

type FederationRecord struct {
	AccountId string `db:"id"`
	MemoType  string `db:"memo_type"`
	Memo      string `db:"memo"`
}

type ReverseFederationRecord struct {
	Name string `db:"name"`
}
