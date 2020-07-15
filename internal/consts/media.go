package consts

const (
	IndicesMediaMapping = `
{
  "settings": {
    "analysis": {
      "analyzer": {
        "pinyin_analyzer": {
          "tokenizer": "my_pinyin"
        }
      },
      "tokenizer": {
        "my_pinyin": {
          "type": "pinyin",
          "keep_separate_first_letter": false,
          "keep_full_pinyin": true,
          "keep_original": true,
          "limit_first_letter_length": 16,
          "lowercase": true,
          "remove_duplicated_term": true
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "@timestamp": {
        "type": "date"
      },
      "@version": {
        "type": "keyword"
      },
      "alias": {
        "type": "keyword",
        "fields": {
          "pinyin": {
            "type": "text",
            "store": false,
            "term_vector": "with_offsets",
            "analyzer": "pinyin_analyzer"
          }
        }
      },
      "casts": {
        "type": "keyword",
        "fields": {
          "pinyin": {
            "type": "text",
            "store": false,
            "term_vector": "with_offsets",
            "analyzer": "pinyin_analyzer"
          }
        }
      },
      "cover": {
        "type": "text",
        "index": false
      },
      "created_at": {
        "type": "date"
      },
      "current_season": {
        "type": "long"
      },
      "directors": {
        "type": "keyword",
        "fields": {
          "pinyin": {
            "type": "text",
            "store": false,
            "term_vector": "with_offsets",
            "analyzer": "pinyin_analyzer"
          }
        }
      },
      "duration": {
        "type": "text",
        "index": false
      },
      "episodes_count": {
        "type": "long"
      },
      "episodes_update": {
        "type": "long"
      },
      "genres": {
        "type": "keyword",
        "fields": {
          "pinyin": {
            "type": "text",
            "store": false,
            "term_vector": "with_offsets",
            "analyzer": "pinyin_analyzer"
          }
        }
      },
      "id": {
        "type": "long"
      },
      "language": {
        "type": "keyword"
      },
      "official_website": {
        "type": "text",
        "index": false
      },
      "original_title": {
        "type": "keyword"
      },
      "rate": {
        "type": "float"
      },
      "region": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "released": {
        "type": "text",
        "index": false
      },
      "released_timestamp": {
        "type": "date"
      },
      "source": {
        "type": "long"
      },
      "subject": {
        "type": "long"
      },
      "subtype": {
        "type": "keyword"
      },
      "summary": {
        "type": "text"
      },
      "title": {
        "type": "keyword",
        "fields": {
          "pinyin": {
            "type": "text",
            "store": false,
            "term_vector": "with_offsets",
            "analyzer": "pinyin_analyzer"
          }
        }
      },
      "torrent_num": {
        "type": "long"
      },
      "updated_at": {
        "type": "date"
      }
    }
  }
}
`
	IndicesMediaConst = "media"
	MediaDefaultPageSize = 20
	SubtypeMovie = "movie"
	SubtypeTV = "tv"
)