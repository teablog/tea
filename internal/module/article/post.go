package article

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"github.com/teablog/tea/internal/helper"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"
)

const PageSize = 10

var (
	Post _post
)

type _post struct{}

type Article struct {
	Title                       string             `json:"title"`
	Keywords                    string             `json:"keywords"`
	Label                       string             `json:"label"`
	Cover                       string             `json:"cover"`
	CoverRaw                    string             `json:"cover_raw"`
	Description                 string             `json:"description"`
	Author                      string             `json:"author"`
	Date                        time.Time          `json:"date"`
	LastEditTime                time.Time          `json:"last_edit_time"`
	Content                     string             `json:"content"`
	Email                       string             `json:"email"`
	Github                      string             `json:"github"`
	Key                         string             `json:"key"`
	Id                          string             `json:"id"`
	Topic                       string             `json:"topic"`
	FilePath                    string             `json:"-"`
	WechatSubscriptionQrcodeRaw string             `json:"wechat_subscription_qrcode_raw"`
	WechatSubscriptionQrcode    string             `json:"wechat_subscription_qrcode"`
	WechatSubscription          string             `json:"wechat_subscription"`
	Md5                         string             `json:"md5"`
	Pv                          int                `json:"pv"`
	Type                        consts.ArticleType `json:"type"`
	Highlight                   []string           `json:"highlight"`
}

func (p *_post) List(ctx *gin.Context, page int) (int, ASlice, error) {
	skip := (page - 1) * PageSize
	var (
		buf bytes.Buffer
	)
	query := map[string]interface{}{
		"from": skip,
		"size": PageSize,
		"sort": map[string]interface{}{
			"last_edit_time": map[string]interface{}{
				"order": "desc",
			},
		},
		"query":   map[string]interface{}{"term": map[string]interface{}{"status": consts.StatusOn}},
		"_source": []string{"author", "title", "description", "topic", "id", "cover", "date", "last_edit_time"},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic(errors.Wrap(err, "json encode 错误"))
	}
	c, data, err := p.Search(buf.String())
	if err != nil {
		return 0, nil, err
	}
	adSensePos := config.Ad.AdSenseFeedsPos()
	list := make(ASlice, 0, len(data)+1)
	for k, v := range data {
		// 插入广告位
		if adSensePos > 0 && k == adSensePos {
			list = append(list, &Article{Type: consts.ArticleTypeAdsense})
		}
		v.Cover = p.ConvertWebp(ctx, v.Cover)
		list = append(list, v)
	}
	return c, list, nil
}

// Get 根据Id获取文章
func (*_post) Get(id string) (bool, *Article, error) {
	resp, err := db.ES.Get(
		consts.IndicesArticleCost,
		id,
	)
	if err != nil {
		return false, nil, errors.Wrapf(err, "get article %s err", id)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, nil, errors.Wrapf(err, "get article %s read response body err", id)
	}
	if resp.IsError() {
		// 文章不存在
		if resp.StatusCode == http.StatusNotFound {
			return false, nil, nil
		}
		return false, nil, errors.Errorf("get article es response %d %s", resp.StatusCode, string(body))
	}
	var r db.ESItemResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return false, nil, errors.Wrapf(err, "get article %s json decode err", id)
	}
	a := new(Article)
	if err := json.Unmarshal(r.Source, a); err != nil {
		return false, nil, errors.Wrapf(err, "get article %s json decode err", id)
	}
	return true, a, nil
}

func (*_post) Search(body string) (int, ASlice, error) {
	resp, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesArticleCost),
		db.ES.Search.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	var eslist db.ESListResponse
	if err := json.NewDecoder(resp.Body).Decode(&eslist); err != nil {
		return 0, nil, err
	}
	type source struct {
		Source *Article `json:"_source"`
	}
	hits := make([]*source, 0)
	if err := json.Unmarshal(eslist.Hits.Hits, &hits); err != nil {
		return 0, nil, err
	}
	m := make(ASlice, 0)
	for _, v := range hits {
		m = append(m, v.Source)
	}
	return eslist.Hits.Total.Value, m, nil
}

func (p *_post) All(source []string) (ASlice, error) {
	body := fmt.Sprintf(`{ "_source": ["%s"], "size": 10000, "query": { "term": {"status": %d}} }`, strings.Join(source, `","`), consts.StatusOn)
	_, data, err := p.Search(body)
	return data, err
}

func (p *_post) Today(source []string) (ASlice, error) {
	//body := fmt.Sprintf(`{ "_source": ["%s"], "size": 10000, "query": { "bool": { "filter": { "range": {  }, "term": {"status": %d}  } } } }`, strings.Join(source, `","`), consts.StatusOn)
	body := fmt.Sprintf(`{
  "_source": ["%s"],
  "size": 200,
  "query": {
    "bool": {
      "must": [
        {
          "range": {
            "last_edit_time": {
              "gte": "%s"
            }
          }
        },
        {
          "term": {
            "status": %d
          }
        }
      ]
    }
  }
}`, strings.Join(source, ","), time.Now().Format(consts.EsTimeYMDFormat), consts.StatusOn)
	_, data, err := p.Search(body)
	return data, err
}

func (p *_post) Labels(size int) (ASlice, error) {
	var (
		buf bytes.Buffer
	)
	query := map[string]interface{}{
		"_source": []string{"label", "id"},
		"size":    size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": map[string]interface{}{
					"script": map[string]interface{}{
						"script": map[string]interface{}{
							"source": "doc['label'].size() > 0",
						},
					},
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic(errors.Wrap(err, "json encode错误"))
	}
	_, data, err := p.Search(buf.String())
	return data, err
}

// 刷新文档
func (*_post) Flush() error {
	resp, err := db.ES.Indices.Flush(
		func(request *esapi.IndicesFlushRequest) {
			request.Index = []string{consts.IndicesArticleCost}
		},
	)
	if err != nil {
		return errors.Wrapf(err, "flush %s err", consts.IndicesArticleCost)
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return errors.Wrapf(err, "flush %s err", consts.IndicesArticleCost)
	}
	return nil
}

func (*_post) DeleteByIds(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	query := fmt.Sprintf(`{
  "query": {
    "terms": {
      "id": ["%s"]
    }
  },
  "script": {
    "source": "ctx._source.status=3",
    "lang": "painless"
  }
}`, strings.Join(ids, `","`))
	resp, err := db.ES.UpdateByQuery(
		[]string{consts.IndicesArticleCost},
		func(request *esapi.UpdateByQueryRequest) {
			request.Body = strings.NewReader(query)
		},
	)
	if err != nil {
		return errors.Wrapf(err, "[es] delete by ids err")
	}
	defer resp.Body.Close()
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "[es] delete by ids response, read body err")
	}
	if resp.IsError() {
		return errors.Errorf("[es] delete by ids body: %s\tresponse %d, %s", query, resp.StatusCode, string(bt))
	}
	type UpdateResp struct {
		Total int
	}
	ur := &UpdateResp{}
	if err := json.Unmarshal(bt, ur); err != nil {
		return errors.Wrapf(err, "[es] update_by_query json decode err")
	}
	if ur.Total != len(ids) {
		return errors.Errorf("[es] delete by ids total not equal: %d, %d", ur.Total, len(ids))
	}
	return nil
}

// ConvertWebp chrome 和 Android 使用webp响应, 其他设别正常返回数据
func (*_post) ConvertWebp(ctx *gin.Context, image string) string {
	ext := path.Ext(image)
	if helper.Image.WebPSupportExt(ext) {
		ua := ctx.Request.UserAgent()
		if strings.Contains(ua, "Chrome") || strings.Contains(ua, "Android") {
			return strings.Replace(image, ext, ".webp", 1)
		}
	}
	return image
}

// ConvertContentWebP 图片生成webp
func (p *_post) ConvertContentWebP(ctx *gin.Context, content string) string {
	matched, err := regexp.MatchString(consts.MarkDownImageRegex, content)
	if err != nil {
		return content
	}
	if matched {
		re, _ := regexp.Compile(consts.MarkDownImageRegex)
		for _, v := range re.FindAllStringSubmatch(content, -1) {
			filename := v[2] + v[3]
			if !strings.HasPrefix(filename, "http") {
				filename = config.Global.CdnHost() + "/" + strings.TrimLeft(filename, "/")
			}
			WebP := Post.ConvertWebp(ctx, filename)
			if WebP != filename {
				// 替换文件image路径
				rebuild := strings.ReplaceAll(v[0], v[2]+v[3], WebP)
				content = strings.ReplaceAll(content, v[0], rebuild)
			}
		}
	}
	return content
}

// 拼接文章id md5(user.key-topic-文件名称)
func (p *_post) GenerateId(topic, key, filename string) string {
	return helper.Md532([]byte(fmt.Sprintf("%s-%s-%s", topic, key, filename)))
}

type ASlice []*Article

func (rows ASlice) MapMd5() map[string]struct{} {
	m := make(map[string]struct{}, len(rows))
	for _, v := range rows {
		m[v.Md5] = struct{}{}
	}
	return m
}

func (rows ASlice) MapId() map[string]struct{} {
	m := make(map[string]struct{}, len(rows))
	for _, v := range rows {
		m[v.Id] = struct{}{}
	}
	return m
}
