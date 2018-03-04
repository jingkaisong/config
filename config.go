package config

import (
	"encoding/json"
	"io/ioutil"
    "time"
	"os"
	"path/filepath"
	"regexp"
)

//DBconf DBconfig
type MySQLConf struct {
	NameServer string `json:"nameserver"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	DBName     string `json:"dbname"`
	User       string `json:"user"`
	Passwd     string `json:"pwd"`
}

type RedisConf struct {
	NameServer string `json:"nameserver"`
}

type PushServerConf struct {
    NameServer string `json:"nameserver"`
}

/*
type MySQLConfs map[string]MySQLConf

type Config struct {
	mySqlConf MySQLConfs
}
*/
type ENVConfig struct {
    Mode string `json:"mode"`
}

var (
    ConfigData map[string][]byte
    LastModTime time.Time
)

func Parse(confFilePath string) (configMap map[string][]byte, err error) {
	file, err := os.Open(filepath.Join(confFilePath, "config.json"))
	if err != nil {
		return
	}
    defer file.Close()
    fileInfo, err := file.Stat()
    if err != nil {
        return
    }
    if LastModTime == fileInfo.ModTime() {
        return ConfigData, nil
    }
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	content = removeAnnotation(content)
	var configs = make(map[string]interface{}, 1024)
	err = json.Unmarshal(content, &configs)
	if err != nil {
		return
	}

	configMap = make(map[string][]byte, len(configs))
	for k, v := range configs {
		configMap[k], err = json.Marshal(v)
	}
    ConfigData = configMap
    LastModTime = fileInfo.ModTime()
	return
}

//删除json中的注释
func removeAnnotation(src []byte) []byte {
	strReg := "(?P<nocomment>'(?:[^\\\\']|\\\\.)*'|\"(?:[^\\\\\"]|\\\\.)*\")|(?P<coment>//[^\n]*|/\\*(.|\n)*?\\*/)"
	reg := regexp.MustCompile(strReg)
	return reg.ReplaceAll(src, []byte("${nocomment}"))
}
