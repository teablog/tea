package consts

const (
	IndicesMessagesConst  = "messages_test"
	IndicesMessageMapping = `
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
      "article_id": {
        "type": "keyword"
      },
      "content": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "date": {
        "type": "date"
      },
      "id": {
        "type": "keyword"
      },
      "sender": {
        "properties": {
          "avatar_url": {
            "type": "text",
            "index": false
          },
          "created_at": {
            "type": "date",
            "index": false
          },
          "email": {
            "type": "keyword",
            "index": false
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
            "type": "keyword",
            "index": false
          },
          "url": {
            "type": "text",
            "index": false
          },
          "ip": {
            "type": "keyword",
            "index": false
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
