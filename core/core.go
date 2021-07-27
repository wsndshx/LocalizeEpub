package core

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/psanford/memfs"
)

var rootFS *memfs.FS
var target string = ""

func init() {
	//新建一个文件系统
	rootFS = memfs.New()
}

// Open 打开输入的epub文件
func Open(Path string) error {
	// 判断文件路径合法性
	_, err := os.Lstat(Path)
	if err != nil {
		return err
	}
	// 判断是否为epub文件
	if path.Ext(Path) != ".epub" {
		return errors.New("所选的文件不是epub文件")
	}
	// 返回状态
	return nil
}

func Unzip(path string) error {
	// 解压目标文件
	miao, err := zip.OpenReader(path)
	defer miao.Close()
	if err != nil {
		return err
	}

	for _, file := range miao.File {
		unzippath := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			err := rootFS.MkdirAll(unzippath, file.Mode())
			if err != nil {
				return err
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		data, _ := ioutil.ReadAll(fileReader)

		rootFS.WriteFile(unzippath, data, 0755)
	}

	// 列出解压后的目录树
	// err = fs.WalkDir(rootFS, ".", func(path string, d fs.DirEntry, err error) error {
	// 	fmt.Println(path)
	// 	return nil
	// })
	// if err != nil {
	// 	panic(err)
	// }
	return nil
}
