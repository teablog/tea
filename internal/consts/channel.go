package consts

const (
	IndicesChannelConst = "channel"
	TypeChannelPrivate  = "private"
	TypeChannelPublic   = "public"

	GlobalChannelId = "douyacun"
	IndicesChannleMap = `
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
      "created_at": {
        "type": "date"
      },
      "creator": {
        "properties": {
          "avatar_url": {
            "type": "text",
            "index": false
          },
          "created_at": {
            "type": "date"
          },
          "email": {
            "type": "keyword"
          },
          "id": {
            "type": "keyword"
          },
          "name": {
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
          "source": {
            "type": "keyword"
          },
          "url": {
            "type": "text"
          },
          "ip": {
            "type": "keyword"
          }
        }
      },
      "members": {
        "properties": {
          "avatar_url": {
            "type": "text",
            "index": false
          },
          "created_at": {
            "type": "date"
          },
          "email": {
            "type": "keyword"
          },
          "id": {
            "type": "keyword"
          },
          "name": {
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
          "source": {
            "type": "keyword"
          },
          "url": {
            "type": "text"
          },
          "ip": {
            "type": "keyword"
          }
        }
      },
      "title": {
        "type": "text",
        "fields": {
          "pinyin": {
            "type": "text",
            "store": false,
            "term_vector": "with_offsets",
            "analyzer": "pinyin_analyzer"
          }
        }
      },
      "type": {
        "type": "keyword"
      }
    }
  }
}`
)
