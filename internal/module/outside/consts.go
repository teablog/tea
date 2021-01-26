package outside

const OutsideIndex = "outside"

const OutsideMapping = `
{
  "mappings": {a
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
      }
    }
  }
}
`

const OutsideReg = `<a(.*)?href=\"https\:\/\/www\.douyacun\.com\"(.*?)>Douyacun<\/a>`