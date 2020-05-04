package controllers

import (
	"dyc/internal/db"
	"dyc/internal/helper"
	"github.com/gin-gonic/gin"
)

var Foo *foo

type foo struct{}

func (*foo) Transaction(ctx *gin.Context) {
	tx := db.Write(ctx).Begin()

	type res struct {
		ID uint64 `json:"id"`
		V  string `json:"v"`
	}
	var result res
	if err := tx.Table("foo").Select("id, v").Where("id = ?", 4).Scan(&result).Error; err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, &result)
	return
}
