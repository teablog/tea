package db

import (
	"encoding/json"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"sync"
)

var ES *elasticsearch.Client
var es_once sync.Once

type ESListResponse struct {
	Took    int  `json:"took"`
	Timeout bool `json:"timed_out"`
	Shards  struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64         `json:"max_score"`
		Hits     json.RawMessage `json:"hits"`
	} `json:"hits"`
}

type ESItemResponse struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type"`
	Id     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

type ESDocResponse struct {
	Index       string          `json:"_index"`
	Type        string          `json:"_type"`
	Id          string          `json:"_id"`
	Version     int             `json:"_version"`
	SeqNo       int             `json:"_seq_no"`
	PrimaryTerm int             `json:"_primary_term"`
	Found       bool            `json:"found"`
	Source      json.RawMessage `json:"_source"`
}

func NewElasticsearch(address []string, user, password string) {
	var err error
	es_once.Do(func() {
		cfg := elasticsearch.Config{
			Addresses:             address,
			Username:              user,
			Password:              password,
			CloudID:               "",
			APIKey:                "",
			RetryOnStatus:         nil,
			DisableRetry:          true,
			EnableRetryOnTimeout:  false,
			MaxRetries:            0,
			DiscoverNodesOnStart:  false,
			DiscoverNodesInterval: 0,
			EnableMetrics:         false,
			EnableDebugLogger:     false,
			RetryBackoff:          nil,
			Transport:             nil,
			Logger:                nil,
			Selector:              nil,
			ConnectionPoolFunc:    nil,
		}
		ES, err = elasticsearch.NewClient(cfg)
		if err != nil {
			panic(err)
		}
	})
}
