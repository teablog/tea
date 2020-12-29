package consts

const (
	IndicesTopicMapping = `{
  "mappings": {
    "properties": {
      "author": {
        "type": "keyword"
      },
      "content": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "cover": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "date": {
        "type": "date"
      },
      "description": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "email": {
        "type": "keyword"
      },
      "github": {
        "type": "keyword",
        "index": false
      },
      "id": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "key": {
        "type": "keyword"
      },
      "keywords": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "label": {
        "type": "text",
        "fielddata": true
      },
      "last_edit_time": {
        "type": "date"
      },
      "title": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        },
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "topic": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "wechat_subscription": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "wechat_subscription_qrcode": {
        "type": "text",
        "index": false
      },
      "md5": {
        "type": "keyword",
        "index": false
      }
    }
  }
}`
	IndicesArticleCost = "articles_v2"
	MarkDownImageRegex = `!\[(.*)\]\((.*)(.png|.gif|.jpg|.jpeg|.webp)(.*)\)`
	MarkDownLocalJump  = `\[.*\]\((\.?\/?(\w+\/?)+\.md)(.*)\)`
)
