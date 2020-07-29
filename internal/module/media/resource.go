package media

import (
	"bytes"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"text/template"
	"time"
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
	Date      time.Time `json:"updated_at"`
	Url       string
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
		data = append(data, Resource.toMap(v.(map[string]interface{})["_source"], ""))
	}
	return
}

func (*_Resource) toMap(data interface{}, author string) interface{} {
	if v, ok := data.(map[string]interface{}); ok {
		var (
			f = map[string]interface{}{
				"id":             v["id"],
				"rate":           int64(v["rate"].(float64)),
				"title":          v["title"],
				"last_edit_time": v["released"],
				"description":    v["summary"],
				"author":         author,
				"cover":          "",
			}
			genres []string
		)
		// author
		if a, ok := v["genres"].([]interface{}); ok && author == "" {
			if len(a) > 2 {
				a = a[:2]
			}
			if len(a) == 0 {
				if a, ok = v["region"].([]interface{}); ok {
					if len(a) > 2 {
						a = a[:2]
					}
				}
			}
			for _, g := range a {
				genres = append(genres, g.(string))
			}
			f["author"] = strings.Join(genres, "/")
		}
		// cover
		if cover, ok := v["cover"]; ok && len(cover.(string)) > 2 {
			f["cover"] = fmt.Sprintf("%s/images/media/%s", consts.Host, cover)
		}
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
	if len(data.Cover) > 2 {
		data.Cover = fmt.Sprintf("%s%s%s", consts.Host, "/images/media/", data.Cover)
	} else {
		data.Cover = ""
	}
	// torrents
	resTorrent, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesMediaTorrentConst),
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
	for _, t := range hitsTorrent.Hits.Hits {
		data.Torrents = append(data.Torrents, t.Source)
	}

	// shares
	var hitsShare _shareResp
	resShare, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesMediaShareConst),
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
	for _, v2 := range hitsShare.Hits.Hits {
		data.Shares = append(data.Shares, v2.Source)
	}
	return
}

func (*_Resource) ToArticle(data _article) (res map[string]interface{}, err error) {
	data.Url = fmt.Sprintf("%s/search/media", consts.Host)

	text := `
![]({{.Cover}})

{{if .Alias}}
**更多中文名：**  {{range $k, $v := .Alias}} {{if $k}}/{{end}} {{$v}} {{end}}
{{end}}

{{if .Language}}
**对白语言：** {{range $k, $v := .Language}} {{if $k}}/{{end}} {{$v}} {{end}}
{{end}}

{{if .Duration}}
**片长：** {{.Duration}}
{{end}}

{{if .Casts}}
**演员：** {{range $k, $v := .Casts}} {{if $k}}/{{end}} <a href='{{html $.Url}}?q=casts:"{{$v}}"' target="_blank">{{$v}}</a> {{end}}
{{end}}

{{if .Released}}
**上映时间：**  {{.Released}}
{{end}}

{{if .Directors}}
**导演：** {{range $k, $v := .Directors}} {{if $k}}/{{end}} <a href='{{html $.Url}}?q=directors:"{{$v}}"' target="_blank">{{$v}}</a>{{end}}
{{end}}

{{if .Genres}}
**类型：** {{range $k, $v := .Genres}} {{if $k}}/{{end}} <a href='{{html $.Url}}?q=genres:"{{$v}}"' target="_blank">{{$v}}</a> {{end}}
{{end}}

{{if .Summary}}
**剧情：**  {{.Summary}}
{{end}}

{{if .Torrents}}
## 下载地址

{{range $k, $v := .Torrents}}
- {{$v.Desc}} <a href="javascript:void(0)" onclick="alert('复制成功')" class="copy" data-clipboard-text="{{$v.Torrent}}">点此复制</a>
{{end}}
{{end}}

{{if .Shares}}
## 在线观看
{{range $k, $v := .Shares}}
- <a href="{{$v.Source}}" target="_blank">{{$v.Desc}} {{$v.Source}}</a>
{{end}}
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

	err = tmpl.Execute(&buf, data)

	res = map[string]interface{}{
		"author":      "",
		"content":     buf.String(),
		"date":        data.Date,
		"description": "",
		"keywords":    data.Title + "," + strings.Join(data.Alias, ",") + ",迅雷下载,种子下载,在线观看",
		"title":       data.Title,
	}
	return
}

func (*_Resource) Search(page int, search string) (total int64, data []interface{}, err error) {
	// 解析查询字段
	value := ""
	if p := strings.Index(search, ":"); p != -1 {
		value = search[p+1:]
	}
	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesMediaConst),
		db.ES.Search.WithQuery(search),
		db.ES.Search.WithFrom((page-1)*consts.MediaDefaultPageSize),
		db.ES.Search.WithSize(consts.MediaDefaultPageSize),
		db.ES.Search.WithSource("id", "title", "region", "genres", "released", "rate", "summary", "cover"),
	)
	if err != nil {
		panic(errors.Wrap(err, "es 查询错误"))
	}
	if res.IsError() {
		panic(errors.Errorf("[%s] ES error", res.Status()))
	}
	defer res.Body.Close()
	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		panic(errors.Wrapf(err, "search/media?%s json decode failed", search))
	}
	total = int64(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	for _, v := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		data = append(data, Resource.toMap(v.(map[string]interface{})["_source"], value))
	}
	return
}
