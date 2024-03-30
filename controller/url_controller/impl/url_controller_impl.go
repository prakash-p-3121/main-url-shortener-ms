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
	var req url_model.ShortenUrlReq
	err := ctx.BindJSON(&req)
	if err != nil {
		badReqErr := errorlib.NewBadReqError("payload-serialization" + err.Error())
		badReqErr.SendRestResponse(ctx)
		return
	}

	resp, appErr := controller.UrlService.ShortenUrl(&req)
	if appErr != nil {
		appErr.SendRestResponse(ctx)
		return
	}
	restlib.OkResponse(ctx, resp)
}

func (controller *UrlControllerImpl) FindLongUrl(restCtx restlib.RestContext) {
	ginRestCtx, ok := restCtx.(*restlib.GinRestContext)
	if !ok {
		internalServerErr := errorlib.NewInternalServerError("Expected GinRestContext")
		internalServerErr.SendRestResponse(ginRestCtx.CtxGet())
		return
	}
	ctx := ginRestCtx.CtxGet()
	shortUrlEncHash := ctx.Param("url-hash")
	if restlib.TrimAndCheckForEmptyString(&shortUrlEncHash) {
		badReqErr := errorlib.NewBadReqError("url-hash-empty")
		badReqErr.SendRestResponse(ctx)
		return
	}
	resp, appErr := controller.UrlService.FindLongUrl(&shortUrlEncHash)
	if appErr != nil {
		appErr.SendRestResponse(ctx)
		return
	}
	restlib.OkResponse(ctx, resp)
}

func (controller *UrlControllerImpl) FindTopDomains(ctx restlib.RestContext) {

}
