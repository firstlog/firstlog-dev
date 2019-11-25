package collect

import (
	"log"
	"os"
	"reflect"
	"sync"
	"syscall"
	"time"
	"strconv"
	"context"

	"firstlog/collect/tail"
	"firstlog/output"
)

var mutex     sync.Mutex
var GoTailRum = 0

func Tail(source string,ctx context.Context) {
	inode      := file2Inode(source)
	device     := file2Device(source)
	offset     := Metadata.get(inode+device)
	createTime := time.Now().Format("2006-01-02 15:04:05")

	Metadata.add(inode+device,source,offset,createTime,inode,device)

	t, err := tail.TailFile(source, tail.Config{
		Follow:    true,
		MustExist: true,
		Poll:      true,
		Location:  &tail.SeekInfo{Offset: offset, Whence: 0},
	})
	if err != nil {
		log.Println(err)
		return
	}

	mutex.Lock()
	GoTailRum ++
	mutex.Unlock()

	tailDone := make(chan bool)
	go func(t *tail.Tail) {
		for line := range t.Lines {
			output.Storage <- line.Text
		}

		Metadata.add(inode+device,source,t.DelOffset,createTime,inode,device) // 删除文件时更新最后的offset到内存
		tailDone <- true
	}(t)

	for {
		select {
		case <-ctx.Done():
			mutex.Lock()
			GoTailRum --
			mutex.Unlock()

			offset, err := t.Tell()
			if err != nil {
				log.Println(err)
			}
			Metadata.add(inode+device,source,offset,createTime,inode,device) // 删除文件时更新最后的offset到内存
			return
		case <-tailDone:
			mutex.Lock()
			GoTailRum --
			mutex.Unlock()

			err := t.Stop()
			if err != nil {
				log.Println(err)
				return
			}
			return
		case <-time.After(time.Second * 1):
			offset,_  := t.Tell()
			Metadata.add(inode+device,source,offset,createTime,inode,device) // 循环更新offset到内存
		}
	}
}

func file2Inode(source string) string {
	fileInfo, err := os.Stat(source)
	if err != nil {
		log.Println(err)
		return ""
	}
	return strconv.FormatUint(reflect.ValueOf(fileInfo.Sys()).Elem().FieldByName("Ino").Uint(), 10)
}

func file2Device(source string) string {
	fd, err := os.Open(source)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer fd.Close()

	var st syscall.Stat_t
	err = syscall.Fstat(int(fd.Fd()),&st)
	if err != nil {
		log.Println(err)
	}
	return strconv.FormatUint(st.Dev,10)
}
