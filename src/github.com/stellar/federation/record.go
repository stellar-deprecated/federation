package federation

type Record struct {
	Username  string `db:"username" json:"stellar_address"`
	AccountId string `db:"account_id" json:"account_id"`
	MemoType  string `db:"memo_type" json:"memo_type,omitempty"`
	Memo      string `db:"memo" json:"memo,omitempty"`
}
