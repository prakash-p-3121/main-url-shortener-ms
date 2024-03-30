package url_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prakash-p-3121/restlib"
)

func ShortenUrl(c *gin.Context) {
	ginRestCtx := restlib.NewGinRestContext(c)
	controller := NewUrlController()
	controller.ShortenUrl(ginRestCtx)
}

func RedirectToLongUrl(c *gin.Context) {
	ginRestCtx := restlib.NewGinRestContext(c)
	controller := NewUrlController()
	controller.RedirectToLongUrl(ginRestCtx)
}

func FindTopDomains(c *gin.Context) {
	ginRestCtx := restlib.NewGinRestContext(c)
	controller := NewUrlController()
	controller.FindTopDomains(ginRestCtx)
}
