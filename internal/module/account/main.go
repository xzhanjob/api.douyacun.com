package account

import (
	"bytes"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"time"
)

type Accouter interface {
	GetName() string
	GetEmail() string
	GetId() string
	GetAvatarUrl() string
	GetUrl() string
	Source() string
}

type Account struct {
	Name      string    `json:"name"`
	Source    string    `json:"source"`
	Id        string    `json:"id"`
	Url       string    `json:"url"`
	AvatarUrl string    `json:"avatar_url"`
	Email     string    `json:"email"`
	CreateAt  time.Time `json:"create_at"`
	Ip        string    `json:"ip"`
}

func NewAccount() *Account {
	return &Account{}
}

func (a *Account) Create(ctx *gin.Context, i Accouter) (data *Account, err error) {
	id := helper.Md516([]byte(i.GetId() + i.Source()))
	var buf bytes.Buffer
	data = &Account{
		Name:      i.GetName(),
		Source:    i.Source(),
		Id:        i.GetId(),
		Url:       i.GetUrl(),
		AvatarUrl: i.GetAvatarUrl(),
		Email:     i.GetEmail(),
		CreateAt:  time.Now(),
		Ip:        helper.RealIP(ctx.Request),
	}
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		panic(errors.Wrap(err, "Account create json encode failed"))
	}
	res, err := db.ES.Index(
		consts.IndicesAccountConst,
		strings.NewReader(buf.String()),
		db.ES.Index.WithDocumentID(id),
	)
	if err != nil {
		panic(errors.Wrap(err, "es Account index create failed"))
	}
	defer res.Body.Close()
	if res.IsError() {
		d, _ := ioutil.ReadAll(res.Body)
		panic(errors.Errorf("[%s] Account create failed, response: %s", res.StatusCode, d))
	}
	data.Id = id
	return
}

func (a *Account) EnableAccess() bool {
	res, err := db.ES.Exists(
		consts.IndicesAccountConst,
		a.Id,
	)
	if err != nil {
		panic(errors.Wrap(err, "account enable access ES exists failed"))
	}
	if res.IsError() {
		return false
	}
	return true
}

func (a *Account) Mget(ids []string) *[]Account {
	// 查询members
	type responseBody struct {
		Docs []struct {
			Source Account `json:"_source"`
			Id     string  `json:"_id"`
		} `json:"docs"`
	}
	idsStr, err := json.Marshal(map[string]interface{}{
		"ids": ids,
	})
	if err != nil {
		panic(errors.Wrap(err, "channel create json encode failed"))
	}
	res, err := db.ES.Mget(
		strings.NewReader(string(idsStr)),
		db.ES.Mget.WithIndex(consts.IndicesAccountConst),
	)
	if err != nil {
		panic(errors.Wrap(err, "es mget request failed"))
	}
	defer res.Body.Close()
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		panic(errors.Errorf("[%d] es mget response: %s", res.StatusCode, resp))
	}
	var r responseBody
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		panic(errors.Wrap(err, "json decode es response body failed"))
	}

	var (
		m []Account
	)
	for _, v := range r.Docs {
		v.Source.Id = v.Id
		m = append(m, v.Source)
	}

	return &m
}

type Cookie struct {
	*Account
	Md5 string `json:"md5"`
}

func (a *Account) SetCookie(ctx *gin.Context) {
	var (
		c   Cookie
		err error
	)
	c.Account = a
	data, err := json.Marshal(a)
	if err != nil {
		panic(errors.Wrap(err, "set cookie json encode failed"))
	}
	c.Md5 = helper.Md532(data)
	cookie, err := json.Marshal(c)
	if err != nil {
		panic(errors.Wrap(err, "set cookie json encode failed"))
	}
	ctx.SetCookie(consts.CookieName, string(cookie), 604800, "/", ".douyacun.com", false, true)
}

func (a *Account) ExpireCookie(ctx *gin.Context) {
	ctx.SetCookie(consts.CookieName, "", -1, "/", ".douyacun.com", false, true)
}

func (c *Cookie) VerifyCookie() bool {
	// 验证cookie完整性
	a, err := json.Marshal(c.Account)
	if err != nil {
		return false
	}
	logger.Debugf("\naccount: %s\nc.md5 %s\nmd5 calculate: %s", string(a), c.Md5, helper.Md532(a))
	if c.Md5 != helper.Md532(a) {
		return false
	}
	return true
}
