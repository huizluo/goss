package objects

import (
	"fmt"
	"github.com/huizluo/goss/app/apiserver/locate"
	"github.com/huizluo/goss/pkg/utils"
	"io"
	"log"
	"net/http"
	"net/url"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.IsExist(hash) {
		return http.StatusOK, nil
	}

	log.Println("obj is not exist ,start to save")
	stream, e := putStream(hash, size)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}

	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	hash, _ = url.PathUnescape(hash)
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}
