package objects

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"goss/app/dataserver/locate"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func sendFile(w io.Writer, file string) {
	f, e := os.Open(file)
	if e!=nil{
		log.Println(e)
		return
	}
	defer f.Close()
	gzipStream,e:=gzip.NewReader(f)
	if e!=nil{
		log.Println(e)
		return
	}
	io.Copy(w, gzipStream)
	gzipStream.Close()
}

func getFile(name string) string {
	files,_:=filepath.Glob(os.Getenv("STORAGE_PATH") + "/objects/" + name + ".*")
	if len(files)!=1{
		return ""
	}
	file:=files[0]
	h:=sha256.New()
	sendFile(h,file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(file, ".")[2]
	if d != hash {
		log.Println("object hash mismatch, remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	return file
}
