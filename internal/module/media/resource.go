package media

import (
	"bytes"
	"dyc/internal/consts"
	"dyc/internal/db"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

var Resource _Resource

type _Resource struct{}

type _article struct {
	ID        int64 `json:"id"`
	Cover     string
	Title     string
	Alias     []string
	Language  []string
	Duration  string
	Casts     []string
	Released  string
	Directors []string
	Genres    []string
	Summary   string
	Torrents  []_torrent
	Shares    []_share
}

type _torrent struct {
	Torrent string `json:"torrent"`
	Desc    string `json:"desc"`
}

type _share struct {
	Source string `json:"source"`
	Desc   string `json:"desc"`
}

func (*_Resource) Index(page int, subtype string) (total int64, data []interface{}, err error) {
	skip := (page - 1) * consts.MediaDefaultPageSize
	var (
		buf bytes.Buffer
		r   map[string]interface{}
	)
	query := map[string]interface{}{
		"from": skip,
		"size": consts.MediaDefaultPageSize,
		"sort": map[string]interface{}{
			"updated_at": map[string]interface{}{
				"order": "desc",
			},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"range": map[string]interface{}{
							"torrent_num": map[string]interface{}{"gte": 0},
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"subtype": subtype,
						},
					},
				},
			},
		},
		"_source": []string{"id", "title", "region", "genres", "released", "rate", "summary", "cover"},
	}
	err = json.NewEncoder(&buf).Encode(query)
	if err != nil {
		panic(errors.Wrap(err, "json encode 错误"))
	}
	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesMediaConst),
		db.ES.Search.WithBody(&buf),
	)
	if err != nil {
		panic(errors.Wrap(err, "es查询错误"))
	}
	defer res.Body.Close()
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		panic(errors.New(string(resp)))
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		panic(errors.Wrap(err, "json decode 错误"))
	}
	total = int64(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	for _, v := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		data = append(data, Resource.toMap(v.(map[string]interface{})["_source"]))
	}
	return
}

func (*_Resource) toMap(data interface{}) interface{} {
	if v, ok := data.(map[string]interface{}); ok {
		var (
			f = map[string]interface{}{
				"id":          v["id"],
				"rate":        int64(v["rate"].(float64)),
				"title":       v["title"],
				"released":    v["released"],
				"description": v["summary"],
				"author":      "",
				"cover":       "",
			}
			genres []string
		)
		// author
		if f, ok := v["genres"].([]interface{}); ok {
			if len(f) > 2 {
				f = f[:2]
			}
			if len(f) == 0 {
				if f, ok = v["region"].([]interface{}); ok {
					if len(f) > 2 {
						f = f[:2]
					}
				}
			}
			for _, g := range f {
				genres = append(genres, g.(string))
			}
		}
		f["author"] = strings.Join(genres, "/")
		// cover
		f["cover"] = "http://www.douyacun.com/images/media/" + v["cover"].(string)
		return f
	} else {
		return data
	}
}

func (*_Resource) View(id string) (data _article, err error) {
	type _resp struct {
		Source _article `json:"_source"`
	}
	type _torrentResp struct {
		Hits struct {
			Hits []struct {
				Source _torrent `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	type _shareResp struct {
		Hits struct {
			Hits []struct {
				Source _share `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	res, err := db.ES.Get(
		consts.IndicesMediaConst,
		id,
	)
	if err != nil {
		panic(errors.Wrap(err, "es查询错误"))
	}
	defer res.Body.Close()
	if res.IsError() {
		panic(errors.Errorf("[%s] ES error document id = %s", res.Status(), id))
	}
	var resp _resp
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		panic(errors.Wrap(err, "media/:id 接口 es response json decode错误"))
	}
	data = resp.Source

	// torrents
	resTorrent, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesMediaShareConst),
		db.ES.Search.WithQuery(fmt.Sprintf("media_id:%s", id)),
	)
	if err != nil {
		panic(errors.Wrap(err, "es查询错误"))
	}
	defer resTorrent.Body.Close()
	if resTorrent.IsError() {
		panic(errors.Errorf("[%s] ES error media_torrent media_id:%s", res.Status(), id))
	}
	var hitsTorrent _torrentResp
	if err = json.NewDecoder(resTorrent.Body).Decode(&hitsTorrent); err != nil {
		panic(errors.Wrap(err, "media/:id media_torrent json decode 错误"))
	}
	for _, v := range hitsTorrent.Hits.Hits {
		data.Torrents = append(data.Torrents, v.Source)
	}

	// shares
	var hitsShare _shareResp
	resShare, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesMediaTorrentConst),
		db.ES.Search.WithQuery(fmt.Sprintf("media_id:%s", id)),
	)
	if err != nil {
		panic(errors.Wrap(err, "es查询错误"))
	}
	defer resShare.Body.Close()
	if resShare.IsError() {
		panic(errors.Errorf("[%s] ES error media_share media_id:%s", res.Status(), id))
	}
	if err = json.NewDecoder(resShare.Body).Decode(&hitsShare); err != nil {
		panic(errors.Wrap(err, "media:id media_share json decode 错误"))
	}
	for _, v := range hitsShare.Hits.Hits {
		data.Shares = append(data.Shares, v.Source)
	}
	return
}

func (*_Resource) toArticle(data _article) (res map[string]interface{}, err error) {

	text := `
![]({{.cover}})


#  亲爱的新年好

**更多中文名:**  {{.alias}}

**对白语言：** {{.language}}

**片长：** {{.duration}}

**演员：** {{range .casts}} []

**上映时间:**  {{.released}}

**导演：** {{range .directors}}

**类型：** {{range .genres}}

**剧情:**  {{.summary}}



## 下载地址

{{range .torrents}}
- <a href="javascript:void(0)" target="_blank">亲爱的新年好.Happy New Year.2019.HD1080P.X264.AAC.Mandarin.CHS.mp4</a>
{{end}}


## 在线观看
{{range .shares}}
- <a href="javascript:void(0)" target="_blank">亲爱的新年好.Happy New Year.2019.HD1080P.X264.AAC.Mandarin.CHS.mp4</a>
{{end}}


> **攀登者：迅雷下载帮助：**
> 1、想要在线观看，请保存到百度网盘中，没有网盘链接或者链接失效的，请搜索：百度云网盘离线下载教程。
> 2、如需下载电影，请先安装迅雷（旋风），然后右键资源链接，选择迅雷（旋风）下载。
> 3、资源名称中含HD为高清，BD为蓝光 ，Mandarin是普通话，Cantonese是粤语，两者都有为双音轨。

`
	tmpl, err := template.New("toArticle").Parse(text)
	if err != nil {
		return
	}
	var buf bytes.Buffer

	tmpl.Execute(&buf, )
}
