package controllers

import (
	"dyc/internal/helper"
	"dyc/internal/module/media"
	"github.com/gin-gonic/gin"
	"strconv"
)

var Media _Media

type _Media struct {
	Movie _Movie
	TV    _TV
}

type _Movie struct{}
type _TV struct{}

func (*_Media) Index(ctx *gin.Context) {
	subtype := ctx.Param("subtype")
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	total, data, err := media.Resource.Index(page, subtype)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, gin.H{"total": total, "data": data})
	return
}

func (*_Media) View(ctx *gin.Context) {
	id := ctx.Param("id")
	data, err := media.Resource.View(id)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, data)
	return
}
