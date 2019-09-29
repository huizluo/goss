package temp

import (
	"encoding/json"
	"goss/app/dataserver/locate"
	"goss/pkg/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}
func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}
func commitTempObject(datFile string, info *tempInfo) {
	fd,e:=os.Open(datFile)
	if e!=nil{
		log.Println(e)
	}
	d:=url.PathEscape(utils.CalculateHash(fd))
	fd.Close()
	if e :=os.Rename(datFile, os.Getenv("STORAGE_PATH")+"/objects/"+info.Name + "." + d);e!=nil{
		log.Println(e)
	}
	locate.Add(info.hash(),info.id())
}

func del(w http.ResponseWriter, r *http.Request) {
	log.Println("删除obj")
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := os.Getenv("STORAGE_PATH") + "/temp/" + uuid
	datFile := infoFile + ".dat"
	os.Remove(infoFile)
	os.Remove(datFile)
}

func patch(w http.ResponseWriter, r *http.Request) {
	log.Println("保存obj")
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_PATH") + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, e = io.Copy(f, r.Body)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	if actual > tempinfo.Size {
		os.Remove(datFile)
		os.Remove(infoFile)
		log.Println("actual size", actual, "exceeds", tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func readFromFile(uuid string) (*tempInfo, error) {
	f, e := os.Open(os.Getenv("STORAGE_PATH") + "/temp/" + uuid)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	var info tempInfo
	json.Unmarshal(b, &info)
	return &info, nil
}
func post(w http.ResponseWriter, r *http.Request) {
	log.Println("上传对象")
	log.Println(os.Getwd())
	uuid,e:=utils.NewUUID()
	if e!=nil{
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := tempInfo{uuid, name, size}
	e = t.writeToFile()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f,e:=os.Create(os.Getenv("STORAGE_PATH") + "/temp/" + t.Uuid + ".dat")
	if e!=nil{
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
	f.Close()

	w.Write([]byte(uuid))
}

func (t *tempInfo) writeToFile() error {
	f, e := os.Create(os.Getenv("STORAGE_PATH") + "/temp/" + t.Uuid)
	if e != nil {
		return e
	}
	defer f.Close()
	b, _ := json.Marshal(t)
	f.Write(b)
	return nil
}

func put(w http.ResponseWriter, r *http.Request) {
	log.Println("修改对象")
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_PATH") + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.Open(datFile)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	if actual != tempinfo.Size {
		os.Remove(datFile)
		log.Println("actual size mismatch, expect", tempinfo.Size, "actual", actual)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f.Close()
	log.Println("临时对象转正")
	commitTempObject(datFile, tempinfo)
}
