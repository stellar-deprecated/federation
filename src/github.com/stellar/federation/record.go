package federation

type FedRecord struct {
	StellarAddress string `json:"stellar_address"`
	AccountId      string `db:"id" json:"account_id"`
	MemoType       string `db:"memo_type" json:"memo_type,omitempty"`
	Memo           string `db:"memo" json:"memo,omitempty"`
}

type RevFedRecord struct {
	Name string `db:"name"`
}
