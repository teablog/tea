package consts

const (
	IndicesGlobalRegion = "global_region"
	IndicesGlobalRegionMapping = `
{
  "mappings": {
    "properties": {
      "code": {
        "type": "keyword"
      },
      "level": {
        "type": "long"
      },
      "name": {
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
      "name_en": {
        "type": "keyword"
      },
      "name_pinyin": {
        "type": "keyword"
      },
      "path": {
        "type": "keyword"
      },
      "pid": {
        "type": "long"
      }
    }
  },
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
          "lowercase": true
        }
      }
    }
  }
}`
)
