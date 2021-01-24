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
	"github.com/teablog/tea/internal/logger"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const PageSize = 10

var (
	Art _article
)

type _article struct{}

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

func (a *_article) List(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
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
		logger.Wrapf(err, "json encode error")
		helper.Fail(ctx, errors.New("服务器出错了～"))
		return
	}
	c, data, err := a.Search(buf.String())
	if err != nil {
		return
	}
	adSensePos := config.Ad.AdSenseFeedsPos()
	list := make(ASlice, 0, len(data)+1)
	for k, v := range data {
		// 插入广告位
		if adSensePos > 0 && k == adSensePos {
			list = append(list, &Article{Type: consts.ArticleTypeAdsense})
		}
		v.Cover = a.ConvertWebp(ctx, v.Cover)
		list = append(list, v)
	}
	helper.Success(ctx, gin.H{"total": c, "data": list})
	return
}

// Get 根据Id获取文章
func (a *_article) Get(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		helper.Fail(ctx, errors.New("你想找的，不存在奥～"))
		return
	}
	resp, err := db.ES.Get(
		consts.IndicesArticleCost,
		id,
	)
	if err != nil {
		logger.Wrapf(err, "get article %s err", id)
		helper.Fail(ctx, errors.New("服务出错了～"))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Wrapf(err, "get article %s read response body err", id)
		helper.Fail(ctx, errors.New("服务出错了～"))
		return
	}
	if resp.IsError() {
		// 文章不存在
		if resp.StatusCode == http.StatusNotFound {
			helper.FailWithCode(ctx, errors.New("阿弥陀佛，经书在西天取经的路上丢了～"), http.StatusNotFound)
			return
		}
	}
	var r db.ESItemResponse
	if err := json.Unmarshal(body, &r); err != nil {
		logger.Wrapf(err, "get article %s json decode err", id)
		helper.Fail(ctx, errors.New("服务出错了～"))
		return
	}
	at := new(Article)
	if err := json.Unmarshal(r.Source, at); err != nil {
		logger.Wrapf(err, "get article %s json decode err", id)
		helper.Fail(ctx, errors.New("服务出错了～"))
		return
	}
	// 封面
	if len(at.Cover) > 0 {
		at.CoverRaw = config.Global.CdnHost() + at.Cover
		at.Cover = config.Global.CdnHost() + a.ConvertWebp(ctx, at.Cover)
	}
	// 内容图片
	at.Content = a.ConvertContentWebP(ctx, at.Content)
	// 微信二维码
	at.WechatSubscriptionQrcodeRaw = config.Global.CdnHost() + at.WechatSubscriptionQrcode
	at.WechatSubscriptionQrcode = config.Global.CdnHost() + a.ConvertWebp(ctx, at.WechatSubscriptionQrcode)
	helper.Success(ctx, at)
	return
}

func (*_article) Search(body string) (int, ASlice, error) {
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

func (a *_article) All(source []string) (ASlice, error) {
	body := fmt.Sprintf(`{ "_source": ["%s"], "size": 10000, "query": { "term": {"status": %d}} }`, strings.Join(source, `","`), consts.StatusOn)
	_, data, err := a.Search(body)
	return data, err
}

func (a *_article) Today(source []string) (ASlice, error) {
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
	_, data, err := a.Search(body)
	return data, err
}

func (a *_article) Labels(ctx *gin.Context) {
	size, err := strconv.Atoi(ctx.Query("size"))
	if err != nil || size <= 0 {
		size = 20
	}
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
		"sort": map[string]interface{}{
			"last_edit_time": map[string]interface{}{
				"order": "desc",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logger.Wrapf(err, "json encode err")
		helper.Fail(ctx, errors.New("服务出错了～"))
		return
	}
	if _, data, err := a.Search(buf.String()); err != nil {
		logger.Wrapf(err, "labels search err")
		helper.Fail(ctx, errors.New("服务出错了～"))
		return
	} else {
		helper.Success(ctx, data)
		return
	}
}

// 刷新文档
func (*_article) Flush() error {
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

func (*_article) DeleteByIds(ids []string) error {
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
func (*_article) ConvertWebp(ctx *gin.Context, image string) string {
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
func (a *_article) ConvertContentWebP(ctx *gin.Context, content string) string {
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
			WebP := Art.ConvertWebp(ctx, filename)
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
func (a *_article) GenerateId(topic, key, filename string) string {
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
