// Code generated by alioth-center/database-columns. DO NOT EDIT.
// Code generated by alioth-center/database-columns. DO NOT EDIT.
// Code generated by alioth-center/database-columns. DO NOT EDIT.

package model

type whisperuserbalanceCols struct {
	ID                  string
	UserID              string
	BalanceChangeAmount string
	BalanceRemaining    string
	Action              string
	Reason              string
	CreatedAt           string
}

var WhisperUserBalanceCols = &whisperuserbalanceCols{
	ID:                  "id",
	UserID:              "user_id",
	BalanceChangeAmount: "balance_change_amount",
	BalanceRemaining:    "balance_remaining",
	Action:              "action",
	Reason:              "reason",
	CreatedAt:           "created_at",
}
