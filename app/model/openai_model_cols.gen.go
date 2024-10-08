// Code generated by alioth-center/database-columns. DO NOT EDIT.
// Code generated by alioth-center/database-columns. DO NOT EDIT.
// Code generated by alioth-center/database-columns. DO NOT EDIT.

package model

type openaimodelCols struct {
	ID              string
	ClientID        string
	Model           string
	Type            string
	MaxTokens       string
	PromptPrice     string
	CompletionPrice string
	RpmLimit        string
	TpmLimit        string
	CreatedAt       string
	UpdatedAt       string
}

var OpenaiModelCols = &openaimodelCols{
	ID:              "id",
	ClientID:        "client_id",
	Model:           "model",
	Type:            "type",
	MaxTokens:       "max_tokens",
	PromptPrice:     "prompt_price",
	CompletionPrice: "completion_price",
	RpmLimit:        "rpm_limit",
	TpmLimit:        "tpm_limit",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}
