package collect

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"firstlog/collect/fsnotify"
)

const startError  = " directory was not found please restart firstlog after creation"

type RecursiveWatcher struct {
	*fsnotify.Watcher
	Files     chan string
	Folders   chan string
	recursive bool
	ignore    string
	match     string
}

func NewRecursiveWatcher(recursive bool,directory,ignore,match string) (*RecursiveWatcher, error) {
	var folders []string
	var files   []string

	if recursive { // Recursive
		folders = SubFoldersRecursive(directory,ignore)
		if len(folders) == 0 {
			return nil,errors.New(startError)
		}

		files = SubFilesRecursive(directory,ignore,match)
		if len(files) == 0 {
			log.Println("WARNING: No file were found.")
		}
	}else if !recursive { // !Recursive
		folders = append(folders,directory)
		if len(folders) == 0 {
			return nil,errors.New(startError)
		}

		files = SubFiles(directory,ignore,match)
		if len(files) == 0 {
			log.Println("WARNING: No file were found.")
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	rw := &RecursiveWatcher{Watcher: watcher,recursive: recursive,ignore:ignore,match: match}

	rw.Files   = make(chan string,len(files))
	rw.Folders = make(chan string, len(folders))

	for _, folder := range folders {
		rw.AddFolder(folder)
	}
	for _, file := range files {
		rw.Files <- file
	}
	return rw, nil
}

func (watcher *RecursiveWatcher) AddFolder(folder string) {
	err := watcher.Add(folder)
	if err != nil {
		log.Println("Error watching: ", folder, err)
	}
}

func (watcher *RecursiveWatcher) Run() {
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					if watcher.recursive { // Recursive
						fi, err := os.Stat(event.Name)
						if err != nil {
							log.Println(err)
						} else if fi.IsDir() {
							if !IgnoreFile(watcher.ignore,filepath.Base(event.Name)) {
								_ = filepath.Walk(event.Name, func(path string, info os.FileInfo, err error) error {
									if info.IsDir() {
										err := watcher.RemoveNew(path)
										if err != nil {
											log.Println(err)
										}
										watcher.AddFolder(path)
										return nil
									} else if !info.IsDir() {
										if !IgnoreFile(watcher.ignore,filepath.Base(path)) {
											if MatchFile(watcher.match,filepath.Base(path)) {
												watcher.Files <- path
											}
										}
									}
									return nil
								})
							}
						} else {
							if !IgnoreFile(watcher.ignore,filepath.Base(event.Name)) {
								if MatchFile(watcher.match,filepath.Base(event.Name)) {
									watcher.Files <- event.Name
								}
							}
						}
					}else if !watcher.recursive { // !Recursive
						fi, err := os.Stat(event.Name)
						if err != nil {
							log.Println(err)
						}else if !fi.IsDir() {
							if !IgnoreFile(watcher.ignore,filepath.Base(event.Name)) {
								if MatchFile(watcher.match,filepath.Base(event.Name)) {
									watcher.Files <- event.Name
								}
							}
						}
					}

				}

				if event.Op&fsnotify.Rename== fsnotify.Rename {
					watcher.Remove(event.Name)
				}

			case err := <-watcher.Errors:
				log.Println(err)
			}
		}
	}()
}


func MatchFile(match,s string) bool {
	m, err := regexp.MatchString(match,s)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

// 返回一个切片,包含所有子目录(递归)
func SubFoldersRecursive(path,ignore string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			if IgnoreFile(ignore,name) && name != "." && name != ".." {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

// 返回一个切片,包含所有子目录下的文件(递归)
func SubFilesRecursive(path,ignore,match string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			name := info.Name()

			// 忽略
			if IgnoreFile(ignore,name) {
				return nil
			}

			// 正则匹配
			if MatchFile(match,name) {
				paths = append(paths, newPath)
			}
		}
		return nil
	})
	return paths
}

// 返回一个切片,包含当前目录下的文件
func SubFiles(path,ignore,match string) (paths []string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if !f.IsDir() {

			// 忽略
			if IgnoreFile(ignore,f.Name()) {
				continue
			}

			if MatchFile(match,f.Name()) {
				paths = append(paths,path+"/"+f.Name())
			}
		}
	}
	return paths
}

// 忽略以.和_开头的文件或目录
func IgnoreFile(match,name string) bool {
	//return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
	// reg := regexp.MustCompile("^.*")
	//return reg.MatchString(name)

	m, err := regexp.MatchString(match,name)
	if err != nil {
		log.Fatal(err)
	}
	return m
}
