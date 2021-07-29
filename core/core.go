package core

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/psanford/memfs"
)

var rootFS *memfs.FS
var LogFS *memfs.FS

// 解压的目录
var tmp string = "."

func init() {
	//新建一个文件系统
	rootFS = memfs.New()
	// 创建log
	LogFS = memfs.New()
}

func Start(model string) error {
	// 假设前面有个输入文件路径的
	err := Open(Input)
	if err != nil {
		return err
	}
	// 解压文件
	err = Unzip(Input)
	if err != nil {
		return err
	}
	// 转换文件
	LogInfo.Println("开始转换文件")
	err = Conversion(model)
	if err != nil {
		return err
	}
	// 压缩回去
	LogInfo.Println("开始输出文件")
	err = Zip(Output)
	if err != nil {
		return err
	}
	LogInfo.Println("文件输出完成")

	// 扔掉没用的文件
	rootFS = nil
	runtime.GC()
	rootFS = memfs.New()
	return nil
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
	err = rootFS.MkdirAll(tmp, 0777)
	if err != nil {
		return err
	}

	for _, file := range miao.File {
		unzippath := filepath.Join(tmp, file.Name)
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

		err = rootFS.WriteFile(unzippath, data, file.Mode())
		if err != nil {
			return err
		}
	}
	return nil
}

// target为输出的目标文件
func Zip(target string) error {
	var err error
	//创建zip包文件
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	//创建zip.Writer
	zw := zip.NewWriter(zipfile)
	defer zw.Close()

	miao, err := rootFS.Open(tmp)
	if err != nil {
		return err
	}
	defer miao.Close()

	info, err := miao.Stat()
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(tmp)
	}

	err = fs.WalkDir(rootFS, ".", func(Path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, _ = d.Info()
		//创建文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(Path, tmp))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		//写入文件头信息
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		//写入文件内容
		file, err := fs.ReadFile(rootFS, Path)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, ioutil.NopCloser(bytes.NewReader(file)))

		return err
	})
	if err != nil {
		return err
	}

	return nil
}
