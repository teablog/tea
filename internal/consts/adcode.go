package consts

const (
	IndicesAdCodeMapping = `
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
      "adcode": {
        "type": "keyword"
      },
      "citycode": {
        "type": "keyword"
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
      }
    }
  }
}
`
	IndicesAdCodeConst = "adcode"
)
