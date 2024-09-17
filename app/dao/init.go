package dao

import (
	"embed"

	"github.com/alioth-center/infrastructure/utils/values"
)

var (
	//go:embed rawsql/*
	rawSqlEmbedding embed.FS

	rawSqlList = map[RawsqlKey]string{}
)

type RawsqlKey string

const (
	RawsqlOpenaiClientGetAvailableClients RawsqlKey = "openai_client.get_available_clients.sql"
	RawsqlOpenaiClientGetClientSecrets    RawsqlKey = "openai_client.get_client_secrets.sql"
	RawsqlOpenaiClientListClients         RawsqlKey = "openai_client.list_clients.sql"
	RawsqlWhisperUserGetUserInfo          RawsqlKey = "whisper_user.get_user_info.sql"
	RawsqlOpenaiClientBalanceStatistics   RawsqlKey = "openai_client_balance.statistics.sql"
)

var rawSqlNames = []RawsqlKey{
	RawsqlOpenaiClientGetAvailableClients,
	RawsqlOpenaiClientGetClientSecrets,
	RawsqlOpenaiClientListClients,
	RawsqlWhisperUserGetUserInfo,
	RawsqlOpenaiClientBalanceStatistics,
}

func LoadRawSqlList(driverName string) {
	// load rawsql list, format is rawsql/{driver}/{sql_file}
	prefix := values.BuildStrings("rawsql/", driverName, "/")
	for _, name := range rawSqlNames {
		content, readErr := rawSqlEmbedding.ReadFile(values.BuildStrings(prefix, string(name)))
		if readErr != nil {
			panic(readErr)
		}

		rawSqlList[name] = string(content)
	}
}
