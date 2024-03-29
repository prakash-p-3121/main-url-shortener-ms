package impl

import (
	"github.com/prakash-p-3121/errorlib"
	"github.com/prakash-p-3121/main-url-shortener-ms/model/url_model"
	"github.com/prakash-p-3121/main-url-shortener-ms/service/url_service"
	"github.com/prakash-p-3121/restlib"
)

type UrlControllerImpl struct {
	UrlService url_service.UrlService
}

func (controller *UrlControllerImpl) ShortenUrl(restCtx restlib.RestContext) {
	ginRestCtx, ok := restCtx.(*restlib.GinRestContext)
	if !ok {
		internalServerErr := errorlib.NewInternalServerError("Expected GinRestContext")
		internalServerErr.SendRestResponse(ginRestCtx.CtxGet())
		return
	}

	ctx := ginRestCtx.CtxGet()
	var req url_model.ShortUrl
	err := ctx.BindJSON(&req)
	if err != nil {
		badReqErr := errorlib.NewBadReqError("payload-serialization" + err.Error())
		badReqErr.SendRestResponse(ctx)
		return
	}

	//controller.UrlService.ShortenUrl(&req)

}
func (controller *UrlControllerImpl) FindLongUrl(ctx restlib.RestContext) {

}
func (controller *UrlControllerImpl) FindTopDomains(ctx restlib.RestContext) {

}
