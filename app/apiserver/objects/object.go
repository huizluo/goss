package objects

import (
	"compress/gzip"
	"fmt"
	"github.com/huizluo/goss/app/apiserver/heartbeat"
	"github.com/huizluo/goss/app/apiserver/locate"
	"github.com/huizluo/goss/pkg/elasticsearch"
	"github.com/huizluo/goss/pkg/rs"
	"github.com/huizluo/goss/pkg/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	size := utils.GetSizeFromHeader(r.Header)
	c, e := storeObject(r.Body, url.PathEscape(hash), size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]

	e = elasticsearch.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func get(w http.ResponseWriter, r *http.Request) {

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var e error
	if len(versionId) > 0 {
		version, e = strconv.Atoi(versionId[0])
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	meta, e := elasticsearch.GetMetadata(name, version)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	obj := url.PathEscape(meta.Hash)
	stream, e := GetStream(obj, meta.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//get offset
	offset := utils.GetOffsetFromHeader(r.Header)
	if offset != 0 {
		stream.Seek(offset, io.SeekCurrent)
		w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, meta.Size-1, meta.Size))
		w.WriteHeader(http.StatusPartialContent)
	}
	//if enabled gzip
	acceptGzip := false
	accept_encoding := r.Header["Accept-Encoding"]
	for i := range accept_encoding {
		if accept_encoding[i] == "gzip" {
			acceptGzip = true
			break
		}
	}
	if acceptGzip {
		w.Header().Set("content-encoding", "gzip")
		w2 := gzip.NewWriter(w)
		if _, e := io.Copy(w2, stream); e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w2.Close()
	} else {
		//write buffer
		if _, e := io.Copy(w, stream); e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	stream.Close()
}

func del(w http.ResponseWriter, r *http.Request) {
	//name := strings.Split(r.URL.EscapedPath(), "/")[2]
	name := r.URL.Query()["name"][0]
	version := r.URL.Query()["version"][0]
	fmt.Println("name is ", name)
	fmt.Println("version is ", version)
	v, _ := strconv.Atoi(version)
	meta, e := elasticsearch.GetMetadata(name, v)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//version, e := elasticsearch.SearchLatestVersion(name)
	//if e != nil {
	//	log.Println(e)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//e = elasticsearch.PutMetadata(name, version.Version+1, 0, "")
	//if e != nil {
	//	log.Println(e)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	elasticsearch.DelMetadata(name, v)
	servers := locate.Locate(meta.Hash)
	for _, addr := range servers {
		req, e := http.NewRequest(http.MethodDelete, "http://"+addr+"/objects/"+meta.Hash, nil)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		c := http.Client{}
		res, e := c.Do(req)
		if e != nil || res == nil {
			log.Println(e)
		} else if res.StatusCode != 200 {
			log.Println(ioutil.ReadAll(res.Body))
		}
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if locate.IsExist(url.PathEscape(hash)) {
		e = elasticsearch.AddVersion(name, hash, size)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	ds := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(ds) != rs.ALL_SHARDS {
		log.Println("cannot find enough dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	stream, e := rs.NewRSResumablePutStream(ds, name, url.PathEscape(hash), size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("location", "/temp/"+url.PathEscape(stream.ToToken()))
	w.WriteHeader(http.StatusCreated)
}
