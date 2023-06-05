package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"summer/utils/types"
)

type FileWriter struct {
	path string
}

type fileConfig struct {
	Path string `json:"path"`
}

func NewFileWriter(configRaw json.RawMessage) (*FileWriter, error) {
	var config fileConfig
	if err := json.Unmarshal(configRaw, &config); err != nil {
		return nil, err
	}
	return &FileWriter{
		path: config.Path,
	}, nil
}

func (fw *FileWriter) Write(p []byte) (n int, err error) {
	var logMap map[string]string
	err = json.Unmarshal(p, &logMap)
	if err != nil {
		return len(p), err
	}
	//设置文件名
	var filename = fmt.Sprintf("%s/logs_%s.log", logMap["level"], types.TimeNow().Format("2006-01-02-15"))
	//判断目录是否存在
	if _, err := os.Stat(fw.path + "/" + logMap["level"]); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.MkdirAll(fw.path+"/"+logMap["level"], os.ModePerm) //0777也可以os.ModePerm
		os.Chmod(fw.path+"/"+logMap["level"], os.ModePerm)
	}

	var fullpath = fmt.Sprintf(fw.path+"/%s", filename)

	var file *os.File
	//判断文件是否存在
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		file, err = os.Create(fullpath)
		if err != nil {
			return len(p), err
		}
	} else {
		file, err = os.OpenFile(fullpath, os.O_RDWR|os.O_APPEND, os.ModePerm)
		if err != nil {
			return len(p), err
		}
	}
	defer file.Close()

	return file.Write(p)

}
