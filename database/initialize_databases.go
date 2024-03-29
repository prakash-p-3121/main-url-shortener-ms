package database

import (
	"database/sql"
	"sync"
)

var shardIDToDatabaseConnectionMap *sync.Map
var singleStoreDatabaseConnection *sql.DB

const (
	ShortUrlsTable               string = "short_urls"
	LongToShortUrlsMappingsTable string = "long_to_short_urls_mappings"
	DomainShorteningCountsTable  string = "domain_shortening_counts"
)

func SetShardConnectionsMap(connectionsMap *sync.Map) {
	shardIDToDatabaseConnectionMap = connectionsMap
}

func GetShardConnectionsMap() *sync.Map {
	return shardIDToDatabaseConnectionMap
}

func SetSingleStoreConnection(databaseConnection *sql.DB) {
	singleStoreDatabaseConnection = databaseConnection
}

func GetSingleStoreConnection() *sql.DB {
	return singleStoreDatabaseConnection
}

func GetShardedTableList() []string {
	return []string{ShortUrlsTable, LongToShortUrlsMappingsTable, DomainShorteningCountsTable}
}
