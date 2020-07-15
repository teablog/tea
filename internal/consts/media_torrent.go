package consts

const (
	IndicesMediaTorrentMapping = `
{
  "mappings": {
    "properties": {
      "@timestamp": {
        "type": "date"
      },
      "@type": {
        "type": "keyword"
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
      "id": {
        "type": "keyword"
      },
      "media_id": {
        "type": "long"
      },
      "torrent": {
        "type": "text",
        "index": false
      }
    }
  }
}
`
	IndicesMediaTorrentConst = "media_torrent"
)
