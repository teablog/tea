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
                "index": true,
				"fielddata": true
            },
            "last_edit_time": {
                "type": "date"
            },
            "title": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
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
            }
        }
    }
}`
	IndicesArticleCost = "articles"
	MarkDownImageRegex = `!\[(.*)\]\((.*)(.png|.gif|.jpg|.jpeg|.webp)(.*)\)`
	MarkDownLocalJump  = `\[.*\]\((\.?\/?(\w+\/?)+\.md)(.*)\)`
)
