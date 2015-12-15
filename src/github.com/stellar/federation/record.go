package federation

type Record struct {
  Username  string `db:"username" json:"stellar_address"`
  AccountId string `db:"account_id" json:"account_id"`
}
