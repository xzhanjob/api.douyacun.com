package account

import (
	"bytes"
	"dyc/internal/consts"
	"dyc/internal/db"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"time"
)

var Account *account

type Accouter interface {
	GetName() string
	GetEmail() string
	GetId() string
	GetAvatarUrl() string
	GetUrl() string
	Source() string
}

type account struct {
	Name      string    `json:"name"`
	Source    string    `json:"source"`
	Id        string    `json:"id"`
	Url       string    `json:"url"`
	AvatarUrl string    `json:"avatar_url"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (*account) Create(a Accouter) (data *account, err error) {
	var buf bytes.Buffer
	data = &account{
		Name:      a.GetName(),
		Source:    a.Source(),
		Id:        a.GetId(),
		Url:       a.GetUrl(),
		AvatarUrl: a.GetAvatarUrl(),
		Email:     a.GetEmail(),
		CreatedAt: time.Now(),
	}
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		panic(errors.Wrap(err, "account create json encode failed"))
	}
	res, err := db.ES.Index(
		consts.IndicesAccountConst,
		strings.NewReader(buf.String()),
		db.ES.Index.WithDocumentID(a.GetId()),
	)
	if err != nil {
		panic(errors.Wrap(err, "es account index create failed"))
	}
	defer res.Body.Close()
	if res.IsError() {
		d, _ := ioutil.ReadAll(res.Body)
		panic(errors.Errorf("[%s] account create failed, response: %s", res.StatusCode, d))
	}
	return
}
