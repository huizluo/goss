package objects

import (
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
	f, _ := os.Open(file)
	defer f.Close()
	io.Copy(w, f)
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
