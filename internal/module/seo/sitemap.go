package seo

import (
	"dyc/internal/config"
	"dyc/internal/module/article"
	"errors"
	"fmt"
	"github.com/douyacun/gositemap"
	"github.com/gin-gonic/gin"
)

var Sitemap sitemap

type sitemap struct{}

func (s *sitemap) Generate(ctx *gin.Context) error {
	articles := article.Search.All([]string{"id", "last_edit_time"})
	if len(*articles) < 0 {
		return errors.New("no articles")
	}
	st := gositemap.NewSiteMap()
	st.SetPretty(true)
	st.SetCompress(false)
	st.SetDefaultHost("https://www.douyacun.com")
	st.SetPublicPath(config.GetKey("path::storage_dir").String())
	host := "https://www.douyacun.com/article/%s"

	url := gositemap.NewUrl()
	url.SetLoc("http://www.douyacun.com/")
	url.SetChangefreq(gositemap.Daily)
	url.SetPriority(1)
	st.AppendUrl(url)

	for _, v := range *articles {
		url := gositemap.NewUrl()
		url.SetLoc(fmt.Sprintf(host, v.Id))
		url.SetLastmod(v.LastEditTime)
		url.SetPriority(0.8)
		url.SetChangefreq(gositemap.Monthly)
		st.AppendUrl(url)
	}
	_, err := st.Storage()
	if err != nil {
		return err
	}
	return nil
}
