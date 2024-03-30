package impl

import (
	"github.com/prakash-p-3121/errorlib"
	"github.com/prakash-p-3121/main-url-shortener-ms/model/url_model"
	"github.com/prakash-p-3121/main-url-shortener-ms/service/url_service"
	"github.com/prakash-p-3121/restlib"
	"net/http"
	"strconv"
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

func (controller *UrlControllerImpl) RedirectToLongUrl(restCtx restlib.RestContext) {
	ginRestCtx, ok := restCtx.(*restlib.GinRestContext)
	if !ok {
		internalServerErr := errorlib.NewInternalServerError("Expected GinRestContext")
		internalServerErr.SendRestResponse(ginRestCtx.CtxGet())
		return
	}
	ctx := ginRestCtx.CtxGet()
	shortUrlEncHash := ctx.Param("encoded-url")
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
	ctx.Redirect(http.StatusFound, resp.LongUrl)
}

func (controller *UrlControllerImpl) FindTopDomains(restCtx restlib.RestContext) {
	ginRestCtx, ok := restCtx.(*restlib.GinRestContext)
	if !ok {
		internalServerErr := errorlib.NewInternalServerError("Expected GinRestContext")
		internalServerErr.SendRestResponse(ginRestCtx.CtxGet())
		return
	}
	ctx := ginRestCtx.CtxGet()
	countStr := ctx.Query("count")
	if restlib.TrimAndCheckForEmptyString(&countStr) {
		badReqErr := errorlib.NewBadReqError("count-empty")
		badReqErr.SendRestResponse(ctx)
		return
	}

	count, err := strconv.ParseUint(countStr, 10, 32)
	if err != nil {
		badReqErr := errorlib.NewBadReqError("count-invalid-integer")
		badReqErr.SendRestResponse(ctx)
		return
	}

	domainList, appErr := controller.UrlService.FindTopDomains(count)
	if appErr != nil {
		appErr.SendRestResponse(ctx)
		return
	}
	restlib.OkResponse(ctx, domainList)
}
