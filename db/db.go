package db

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	PathEnvKey        = "GOTERMCONFIG"
	defaultPath       = "/usr/local/goterm"
	defaultConfigName = "command_config.json"
)

type CommandDetail struct {
	Command  []string `json:"command"`
	Project  string   `json:"project"`
	Category string   `json:"category"`
	Remark   string   `json:"remark"`
	Mode     string   `json:"mode"`
}

func Save(data map[string]CommandDetail) error {
	if data == nil {
		return errors.New("empty config to set")
	}

	var (
		f        *os.File
		err      error
		filePath string
	)

	filePath = GetFilePath()

	if f, err = getFile(filePath, true); err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f) // 创建新的 Writer 对象

	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err = w.Write(content); err != nil {
		return err
	}

	if err = w.Flush(); err != nil {
		return err
	}

	return nil
}

func Load() (map[string]CommandDetail, error) {
	var (
		f        *os.File
		err      error
		filePath string
		result   map[string]CommandDetail
	)

	filePath = GetFilePath()

	if f, err = getFile(filePath, false); err != nil {
		return nil, err
	}
	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if len(body) > 0 {
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return result, nil
}

func GetFilePath() string {
	filePath := os.Getenv(PathEnvKey)
	if filePath == "" {
		filePath = defaultPath
	}
	return filePath
}

func getFile(filePath string, trunc bool) (*os.File, error) {
	var (
		f   *os.File
		err error
	)
	ok, _ := pathExists(filePath)
	if !ok {
		err := os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("failed to create goterm config file")
			return nil, err
		}
	}
	filename := fmt.Sprintf("%s/%s", filePath, defaultConfigName)
	if checkFileIsExist(filename) { // 如果文件存在
		if trunc {
			f, err = os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0666) //打开文件
		} else {
			f, err = os.OpenFile(filename, os.O_RDWR, 0666) //打开文件
		}
	} else {
		f, err = os.Create(filename) // 创建文件
	}
	if err != nil {
		fmt.Println("failed to open config file")
		return nil, err
	}

	return f, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func checkFileIsExist(filePath string) bool {
	var exist = true
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
