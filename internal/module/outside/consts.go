package outside

const OutsideIndex = "outside"

const OutsideMapping = `
{
  "mappings": {
    "properties": {
      "create_at": {
        "type": "date"
      },
      "priority": {
        "type": "long"
      },
      "status": {
        "type": "long"
      },
      "title": {
        "type": "keyword",
        "ignore_above": 256
      },
      "url": {
        "type": "keyword",
        "index": false
      },
      "host": {
        "type": "keyword"
      },
      "email": {
        "type": "keyword"
      },
      "reason": {
        "type": "text"
      }
    }
  }
}
`

// 向后断言: ?!</a>
const OutsideReg = `<a[^>]+?href=["'](https\:\/\/www\.douyacun\.com)["'][^>]*>(.*?)</a>`
