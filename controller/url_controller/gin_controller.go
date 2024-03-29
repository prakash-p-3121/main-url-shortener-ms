package product_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prakash-p-3121/restlib"
)

func ShortenUrl(c *gin.Context) {
	ginRestCtx := restlib.NewGinRestContext(c)
	controller := NewUrlController()
	controller.ShortenUrl(ginRestCtx)
}

func FindLongUrl(c *gin.Context) {
	ginRestCtx := restlib.NewGinRestContext(c)
	controller := NewUrlController()
	controller.FindLongUrl(ginRestCtx)
}
