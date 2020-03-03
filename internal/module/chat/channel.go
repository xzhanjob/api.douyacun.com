package chat

import (
	"bytes"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/derror"
	"dyc/internal/module/account"
	"dyc/internal/validate"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"time"
)

var Channel *channel

type channel struct {
	Creator   *account.Account  `json:"creator"`
	Members   []account.Account `json:"members"`
	Title     string            `json:"title"`
	CreatedAt time.Time         `json:"created_at"`
	Id        string            `json:"id"`
}

func (*channel) Create(ctx *gin.Context, v *validate.ChannelCreateValidator) (c map[string]interface{}, err error) {
	if a, ok := ctx.Get("account"); ok {
		type esResponse struct {
			Id string `json:"_id"`
		}
		// 创建channel
		var (
			buf bytes.Buffer
		)
		m := account.NewAccount().Mget(v.Members)
		c = map[string]interface{}{
			"title":      v.Title,
			"creator":    a,
			"created_at": time.Now(),
			"type":       v.Type,
			"members":    *m,
		}
		if err = json.NewEncoder(&buf).Encode(c); err != nil {
			panic(errors.Wrap(err, "channel map json encode failed"))
		}
		res, err := db.ES.Index(
			consts.IndicesChannelConst,
			strings.NewReader(buf.String()),
		)
		if err != nil {
			panic(errors.Wrap(err, "es create channel failed"))
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		if res.IsError() {
			panic(errors.Errorf("[%s] es create channel response: %s", res.StatusCode, body))
		}
		var r esResponse
		if err = json.Unmarshal(body, &r); err != nil {
			panic(errors.Wrapf(err, "response body json decode failed"))
		}
		c["id"] = r.Id
	} else {
		panic(derror.Unauthorized{})
	}
	return
}

func (*channel) Exists() {

}