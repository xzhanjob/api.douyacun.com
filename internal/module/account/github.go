package account

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type _github struct {
	t struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	u struct {
		Id        int64    `json:"id"`
		Url       string `json:"html_url"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatar_url"`
	}
}

func NewGithub() *_github {
	return &_github{}
}

func (g *_github) Token(code string) (err error) {
	params := &gin.H{
		"client_id":     "25fc4f51f48cb5d52edf",
		"client_secret": "e95e1cca0535a3bf9b8a26ed3920b5fbecca873c",
		"code":          code,
	}
	requestBody, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	retries := 3
	var resp *http.Response
	for retries > 0  {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		retries--
	}
	if err != nil {
		panic(errors.Wrap(err, "request github oauth failed"))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		d, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("github oauth/access_token [%s] response: %s", resp.StatusCode, d)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &g.t); err != nil {
		panic(errors.Wrapf(err, "github oauth/access_token json decode failed, response: %s", data))
	}
	if g.t.AccessToken == "" {
		return errors.Errorf("github授权登录失败")
	}
	return
}

func (g *_github) User() (err error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", bytes.NewBuffer([]byte{}))
	authorization := bytes.NewBufferString(g.t.TokenType)
	authorization.WriteString(" ")
	authorization.WriteString(g.t.AccessToken)
	req.Header.Set("Authorization", authorization.String())
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	retries := 3
	var resp *http.Response
	for retries > 0 {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		retries--
	}
	if err != nil {
		panic(errors.Wrap(err, "request github user failed"))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		d, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("github oauth/access_token [%s] response: %s", resp.StatusCode, d)
	}
	if err = json.NewDecoder(resp.Body).Decode(&g.u); err != nil {
		panic(errors.Wrap(err, "github user json decode failed"))
	}
	return
}

func (g *_github) GetName() string {
	return g.u.Name
}

func (g *_github) GetId() string {
	return strconv.FormatInt(g.u.Id, 10)
}

func (g *_github) GetUrl() string {
	return g.u.Url
}

func (g *_github) GetEmail() string {
	return g.u.Email
}

func (g *_github) GetAvatarUrl() string {
	return g.u.AvatarUrl
}

func (g *_github) Source() string {
	return "github"
}
