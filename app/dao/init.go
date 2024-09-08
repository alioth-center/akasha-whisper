package dao

import (
	"embed"
	"github.com/alioth-center/infrastructure/utils/values"
)

var (
	//go:embed rawsql/*.sql
	rawSqlEmbedding embed.FS

	rawSqlList = map[RawsqlKey]string{}
)

type RawsqlKey string

const (
	RawsqlOpenaiClientGetAvailableClients RawsqlKey = "openai_client.get_available_clients.sql"
	RawsqlOpenaiClientGetClientSecrets    RawsqlKey = "openai_client.get_client_secrets.sql"
	RawsqlOpenaiClientListClients         RawsqlKey = "openai_client.list_clients.sql"
	RawsqlWhisperUserGetUserInfo          RawsqlKey = "whisper_user.get_user_info.sql"
)

var (
	rawSqlNames = []RawsqlKey{
		RawsqlOpenaiClientGetAvailableClients,
		RawsqlOpenaiClientGetClientSecrets,
		RawsqlOpenaiClientListClients,
		RawsqlWhisperUserGetUserInfo,
	}
)

func init() {
	// load rawsql list
	for _, name := range rawSqlNames {
		content, readErr := rawSqlEmbedding.ReadFile(values.BuildStrings("rawsql/", string(name)))
		if readErr != nil {
			panic(readErr)
		}

		rawSqlList[name] = string(content)
	}
}
