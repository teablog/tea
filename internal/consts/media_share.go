package consts

const (
	IndicesMediaShareMapping = `
{
  "mappings": {
    "properties": {
      "@timestamp": {
        "type": "date"
      },
      "@version": {
        "type": "text"
      },
      "created_at": {
        "type": "date"
      },
      "desc": {
        "type": "text",
        "index": false
      },
      "episode": {
        "type": "long"
      },
      "id": {
        "type": "keyword"
      },
      "media_id": {
        "type": "long"
      },
      "share_id": {
        "type": "long"
      },
      "share_type": {
        "type": "keyword"
      },
      "source": {
        "type": "text",
        "index": false
      },
      "updated_at": {
        "type": "date"
      }
    }
  }
}`
	IndicesMediaShareConst = "media_share"
)
