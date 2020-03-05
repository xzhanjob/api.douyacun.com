package chat

import (
	"bytes"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/derror"
	"dyc/internal/module/account"
	"dyc/internal/validate"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

var (
	Channel        *channel
	ChannelMembers *channelMembers
)

type channel struct {
	Creator   *account.Account  `json:"creator"`
	Members   []account.Account `json:"members"`
	Title     string            `json:"title"`
	CreatedAt time.Time         `json:"created_at"`
	Id        string            `json:"id"`
	Type      string            `json:"type"`
}

type channelMembers struct {
	m  map[string]*channel
	mu sync.RWMutex
}

type esResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source channel `json:"_source"`
			Id     string  `json:"_id"`
		} `json:"hits"`
	} `json:"hits"`
}

// 创建一个新的channel
func (*channel) Create(ctx *gin.Context, v *validate.ChannelCreateValidator) (c *channel, err error) {
	if a, ok := ctx.Get("account"); ok {
		type esResponse struct {
			Id string `json:"_id"`
		}
		// 创建channel
		var (
			buf   bytes.Buffer
			title = v.Title
		)
		// 标题设置
		m := account.NewAccount().Mget(v.Members)
		if len(*m) == 0 {
			return c, errors.Errorf("members not found")
		}
		if strings.TrimSpace(v.Title) == "" {
			if v.Type == consts.TypeChannelPrivate {
				title = (*m)[0].Name
			} else if v.Type == consts.TypeChannelPublic {
				name := make([]string, 0)
				for _, v := range *m {
					name = append(name, v.Name)
				}
				if len(name) > 20 {
					name = name[:20]
				}
				title = strings.Join(name, "、")
			}
		}
		query := map[string]interface{}{
			"title":      title,
			"creator":    a,
			"created_at": time.Now(),
			"type":       v.Type,
			"members":    *m,
		}
		if err = json.NewEncoder(&buf).Encode(query); err != nil {
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
		c = &channel{
			Id:        r.Id,
			Members:   *m,
			Title:     title,
			Creator:   a.(*account.Account),
			CreatedAt: time.Now(),
		}
	} else {
		panic(derror.Unauthorized{})
	}
	return
}

// 获取channel
func (*channel) Get(ctx *gin.Context, v *validate.ChannelCreateValidator) (c *channel, ok bool) {
	if a, ok := ctx.Get("account"); ok {
		query := fmt.Sprintf(`
{
  "query": {
    "bool": {
      "must": [
        {
          "term": {
            "creator.id": "%s"
          }
        },
        {
          "term": {
            "members.id": "%s"
          }
        },
        {
          "term": {
            "type": "%s"
          }
        }
      ]
    }
  }
}`, a.(*account.Account).Id, v.Members[0], v.Type)
		res, err := db.ES.Search(
			db.ES.Search.WithIndex(consts.IndicesChannelConst),
			db.ES.Search.WithBody(strings.NewReader(query)),
		)
		if err != nil {
			panic(errors.Wrap(err, "channel exists search failed"))
		}
		defer res.Body.Close()
		respBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(errors.Wrapf(err, "[%d] es response body read failed", res.StatusCode))
		}
		if res.IsError() {
			panic(errors.Wrapf(err, "[%d] es response: %s", res.StatusCode, respBody))
		}

		var r esResponse
		if err = json.Unmarshal(respBody, &r); err != nil {
			panic(errors.Wrapf(err, "channel exists es response: %s", respBody))
		}
		if r.Hits.Total.Value > 0 {
			c = &r.Hits.Hits[0].Source
			return c, true
		} else {
			return c, false
		}
	} else {
		panic(derror.Unauthorized{})
	}
}

// channel列表
func (*channel) List(ctx *gin.Context) (*[]channel, error) {
	a, _ := ctx.Get("account")
	query := fmt.Sprintf(`
{
  "query": {
    "bool": {
      "should": [
        {"term": { "creator.id": "%s"}},
        {"term": { "members.id": "%s"}},
		{"term": { "_id": { "value": "douyacun", "boost": 10 }}}
      ]
    }
  }
}
`, a.(*account.Account).Id, a.(*account.Account).Id)
	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesChannelConst),
		db.ES.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		panic(errors.Wrap(err, "es channel index error"))
	}
	defer res.Body.Close()
	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(errors.Wrap(err, "es response body read error"))
	}
	if res.IsError() {
		panic(errors.Errorf("[%d] es response: %s", res.StatusCode, resp))
	}
	var r esResponse
	if err = json.Unmarshal(resp, &r); err != nil {
		panic(errors.Wrapf(err, "json decode failed response: %s", resp))
	}
	var c []channel
	for _, v := range r.Hits.Hits {
		v.Source.Id = v.Id
		c = append(c, v.Source)
	}
	return &c, nil
}

func (m *channelMembers) Join() {

}
