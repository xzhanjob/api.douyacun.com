package elstaticsearch

import (
	"dyc/internal/db"
	"dyc/internal/logger"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

func Create(c *gin.Context) {
	resp, err := db.ES.Index(
		"foo",
		strings.NewReader(`{"title": "foo"}`),
		db.ES.Index.WithDocumentID("1"),
		db.ES.Index.WithRefresh("true"),
	)
	if err != nil {
		logger.Error("insert failed: %s", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("ioutil readall failed: %s", err)
	}
	c.JSON(http.StatusOK, string(data))
}
