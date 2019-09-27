package objects

import (
	"fmt"
	"goss/app/apiserver/locate"
	"goss/pkg/utils"
	"io"
	"net/http"
	"net/url"
)

func storeObject(r io.Reader, hash string,size int64) (int, error) {
	if locate.IsExist(url.PathEscape(hash)){
		return http.StatusOK,nil
	}
	stream, e := putStream(url.PathEscape(hash),size)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}

	reader:=io.TeeReader(r,stream)
	d:=utils.CalculateHash(reader)
	if d!=hash{
		stream.Commit(false)
		return http.StatusBadRequest,fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}