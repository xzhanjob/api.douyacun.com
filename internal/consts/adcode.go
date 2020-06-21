package consts

const (
	IndicesAdCodeMapping = `
{
    "mappings" : {
      "properties" : {
        "adcode" : {
          "type" : "long"
        },
        "citycode" : {
          "type" : "long"
        },
        "name" : {
          "type" : "text",
          "fields" : {
            "keyword" : {
              "type" : "keyword",
              "ignore_above" : 256
            }
          }
        }
      }
    }
  }
`
	IndicesAdCodeConst = "adcode"
)
