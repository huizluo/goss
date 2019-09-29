package main

import (
	"go-oss/src/lib/es"
	"goss/pkg/elasticsearch"
	"log"
)

const MIN_VERSION_COUNT = 5

func main() {
	buckets, e := elasticsearch.SearchVersionStatus(MIN_VERSION_COUNT + 1)
	if e != nil {
		log.Println(e)
		return
	}
	for i := range buckets {
		bucket := buckets[i]
		for v := 0; v < bucket.Doc_count-MIN_VERSION_COUNT; v++ {
			es.DelMetadata(bucket.Key, v+int(bucket.Min_version.Value))
		}
	}
}
