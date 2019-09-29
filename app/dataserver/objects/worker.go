package objects

import (
	"goss/app/dataserver/locate"
	"goss/pkg/utils"
	"io"
	"log"
	"net/url"
	"os"
)

func sendFile(w io.Writer, file string) {
	f, _ := os.Open(file)
	defer f.Close()
	io.Copy(w, f)
}

func getFile(hash string) string {
	file := os.Getenv("STORAGE_PATH") + "/objects/" + hash
	f, _ := os.Open(file)
	d := url.PathEscape(utils.CalculateHash(f))
	f.Close()
	if d != hash {
		log.Println("object hash mismatch, remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	return file
}
