package consts

const (
	SpiderIndices = "spider_%s"
	SpiderIndexTemplate = `{
  "index_patterns": ["spider_*"],  
  "mappings": {
    "properties": {
      "date": {
        "type": "date"
      },
      "spider": {
        "type": "keyword"
      }
    }
  }
}`
)




