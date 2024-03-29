package impl

import "sync"

type UrlRepositoryImpl struct {
	ShardConnectionsMap *sync.Map
}
