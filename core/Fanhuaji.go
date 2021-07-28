package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path"
)

var Mode map[string]string = map[string]string{
	"简体化":   "Simplified",
	"繁體化":   "Traditional",
	"中国化":   "China",
	"香港化":   "Hongkong",
	"台灣化":   "Taiwan",
	"维基简体化": "WikiSimplified",
	"維基繁體化": "WikiTraditional",
}

type Resp struct {
	Code int    `json:"code"`
	Data data   `json:"data"`
	Msg  string `json:"msg"`
}

type data struct {
	Text string `json:"text"`
}

func Conversion(model string) error {
	//转换文件
	data := make(map[string]interface{})
	data["converter"] = Mode[model]
	err := fs.WalkDir(rootFS, ".", func(Path string, d fs.DirEntry, err error) error {
		if path.Ext(Path) == ".xhtml" || path.Ext(Path) == ".ncx" {
			fmt.Println(Path)
			// 读取文件内容
			content, err := fs.ReadFile(rootFS, Path)
			if err != nil {
				panic(err)
			}
			// 生成请求
			data["text"] = string(content)
			bytesData, _ := json.Marshal(data)
			resp, err := http.Post("http://api.zhconvert.org/convert", "application/json", bytes.NewReader(bytesData))
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// 解析返回内容
			body, _ := ioutil.ReadAll(resp.Body)
			var Data Resp
			err = json.Unmarshal(body, &Data)
			if err != nil {
				return err
			}
			if Data.Code == 0 {
				// 回写文件
				rootFS.WriteFile(Path, []byte(Data.Data.Text), d.Type())
			}
		}
		return nil
	})
	return err
}
