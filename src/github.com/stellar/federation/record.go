package federation

type FedRecord struct {
	StellarAddress string `json:"stellar_address"`
	AccountId      string `db:"id" json:"account_id"`
	MemoType       string `db:"type" json:"memo_type"`
	Memo           string `db:"memo" json:"memo"`
}
