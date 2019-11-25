package collect

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"firstlog/initialize"
)

type data struct {
	Source	string    `json:"source"`
	Offset	int64     `json:"offset"`
	Time	string    `json:"time"`
	Inode	string    `json:"inode"`
	Device	string    `json:"device"`
}

type MetadataStruct struct {
	metadataParent map[string]map[string]data
	metadataChild  map[string]data
	path           string
	sync.Mutex
	data
}

func (c *MetadataStruct) add(key,source string,offset int64,time string,inode,device string) {
	c.Source  = source
	c.Offset  = offset
	c.Time	  = time
	c.Inode   = inode
	c.Device  = device

	data1 := c.data

	c.Lock()
	c.metadataChild[key] = data1
	c.Unlock()
}

func (c *MetadataStruct) get(key string) int64 {
	c.Lock()
	value, ok := c.metadataChild[key]
	c.Unlock()
	if ok {
		return value.Offset
	}
	return 0
}

func (c *MetadataStruct) Del(dir map[string]initialize.Task)  {
	existPath := make(map[string]string)

	for k,v := range dir {

		// 递归
		if v.Recursive {
			_ = filepath.Walk(k, func(path string, info os.FileInfo, err error) error {
				file, err := os.Stat(path)
				if err != nil {
					log.Println(err)
					return nil
				}
				if !file.IsDir() {

					if IgnoreFile(v.Ignore, file.Name()) {
						return nil
					}

					if MatchFile(v.Match, file.Name()) {
						existPath[file2Inode(path)+file2Device(path)] = path
					}
				}
				return nil
			})
		}

		// 非递归
		if !v.Recursive {
			files, err := ioutil.ReadDir(k)
			if err != nil {
				log.Println(err)
				continue
			}
			for _, file := range files {
				if !file.IsDir() {

					if IgnoreFile(v.Ignore,file.Name()) {
						continue
					}

					if MatchFile(v.Match,file.Name()) {
						existPath[file2Inode(k+"/"+file.Name())+file2Device(k+"/"+file.Name())] = k+"/"+file.Name()
					}
				}
			}
		}
	}

	c.Lock()
	for k,_ := range c.metadataChild {
		_, ok := existPath[k]
		if !ok {
			delete(c.metadataChild,k)
		}
	}
	c.Unlock()
}

func (c *MetadataStruct) Save2Registry() {
	c.Lock()
	data, err := json.Marshal(&c.metadataChild)
	c.Unlock()

	if err != nil {
		log.Println(err)
		return
	}

	_, err = os.Stat("data")
	if err != nil {
		err := os.Mkdir("data",os.ModePerm)
		if err!=nil{
			log.Println(err)
			return
		}
	}
	err = ioutil.WriteFile("data/registry", data, 0666)
	if err != nil {
		log.Println(err)
		return
	}
}

func NewMetadata() *MetadataStruct {
	res := &MetadataStruct{
		metadataParent: make(map[string]map[string]data),
		metadataChild:  make(map[string]data),
	}
	return res
}