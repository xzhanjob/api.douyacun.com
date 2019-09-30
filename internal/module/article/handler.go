package article

import (
	"context"
	"dyc/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func ListHandler(c *gin.Context) {
	searchResult, err := db.ES.Search().Index(TopicCost).From(0).Size(10).Do(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, "not found")
	}
	var (
		data Article
		res  = make([]Article, 0, 10)
	)
	for _, item := range searchResult.Each(reflect.TypeOf(data)) {
		res = append(res, item.(Article))
	}
	c.JSON(http.StatusOK, gin.H{"total": searchResult.Hits.TotalHits.Value, "data": res})
}
