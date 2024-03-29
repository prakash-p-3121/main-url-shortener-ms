package cfg

import (
	"github.com/prakash-p-3121/restlib"
	restlib_model "github.com/prakash-p-3121/restlib/model"
	"sync"
)

var msConnectionsMap *sync.Map

func SetMsConnectionsMap(connectionsMap *sync.Map) {
	msConnectionsMap = connectionsMap
}

func GetMsConnectionCfg(msName string) (restlib_model.MsConnectionCfg, error) {
	return restlib.RetrieveMsConnectionCfg(msConnectionsMap, msName)
}
