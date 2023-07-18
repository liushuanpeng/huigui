package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	ConfigMap *TotalConf
}

var Conf *Config

func (t *Config) LoadConfig(filePath string) error {
	defer func() error {
		if err := recover(); err != nil {
			return errors.New(fmt.Sprintf("%v", err))
		}
		return nil
	}()
	fd, openErr := os.Open(filePath)
	if openErr != nil {
		return openErr
	}
	defer fd.Close()
	content, readErr := ioutil.ReadAll(fd)
	if readErr != nil {
		return readErr
	}

	configMap := new(TotalConf)
	if err := json.Unmarshal(content, configMap); err != nil {
		return err
	}
	Conf.ConfigMap = configMap
	return nil
}

type HostInfo struct {
	Host string `json:"Host"`
	Port string `json:"Port"`
}

type MysqlConf struct {
	Host         string
	Port         int
	User         string
	Password     string
	DbName       string
	ConnTimeout  int
	ReadTimeout  int
	WriteTimeout int
	MaxOpenConn  int
	MaxIdleConn  int
	Er           string
	RetryCount   int
	TimeInterval int
	DelayTime    int
}

type ImgHost struct {
	host string
	port string
}

type TotalConf struct {
	SaasOrgDB  MysqlConf
	ReqImgHost ImgHost
}
