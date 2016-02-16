package federation

type Response struct {
	StellarAddress string `json:"stellar_address"`
	AccountId      string `json:"account_id"`
	MemoType       string `json:"memo_type,omitempty"`
	Memo           string `json:"memo,omitempty"`
}
