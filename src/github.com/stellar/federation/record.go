package federation

type FedRecord struct {
	StellarAddress string `json:"stellar_address"`
	AccountId      string `db:"id" json:"account_id"`
	MemoType       string `db:"memo_type" json:"memo_type"`
	Memo           string `db:"memo" json:"memo"`
}

type RevFedRecord struct {
	Name string `db:"name"`
}
