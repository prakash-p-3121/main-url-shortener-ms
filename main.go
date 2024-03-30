package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prakash-p-3121/main-url-shortener-ms/cfg"
	"github.com/prakash-p-3121/main-url-shortener-ms/controller/url_controller"

	"github.com/prakash-p-3121/main-url-shortener-ms/database"
	"github.com/prakash-p-3121/mysqllib"
	"github.com/prakash-p-3121/restlib"
	"log"
)

func main() {

	msConnectionsMap, err := restlib.CreateMsConnectionCfg("conf/microservice.toml")
	if err != nil {
		panic(err)
	}
	cfg.SetMsConnectionsMap(msConnectionsMap)

	cfg, err := cfg.GetMsConnectionCfg("database-clustermgt-ms")
	if err != nil {
		panic(err)
	}
	connectionsMap, err := mysqllib.CreateShardConnectionsWithRetry(database.GetShardedTableList(), cfg.Host, cfg.Port)
	if err != nil {
		panic(err)
	}
	log.Println(connectionsMap)
	database.SetShardConnectionsMap(connectionsMap)

	db, err := mysqllib.CreateDatabaseConnectionWithRetryByCfg("conf/database.toml")
	if err != nil {
		panic(err)
	}
	database.SetSingleStoreConnection(db)

	router := gin.Default()
	router.GET("/:encoded-url", url_controller.FindLongUrl)

	routerGroup := router.Group("/public")
	routerGroup.POST("/v1/shorten_url", url_controller.ShortenUrl)

	err = router.Run("127.0.0.1:3000")
	if err != nil {
		panic("Error Starting main url shortener ms")
	}
}
