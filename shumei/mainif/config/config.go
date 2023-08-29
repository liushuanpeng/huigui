package config

import (
	"encoding/json"
	"errors"
	"fmt"
    "net"
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
	Host string
	Port string
}

type KitexConf struct {
    MaxConnections int
    MaxQPS         int
    Weight         int    
}

type BasicConf struct {
    ZkServers             string
    ZkPath                string
    Pwdencrypted          bool
    TextMaxLength         int
    MaxCpuPercent         int32
    MaxMemPercent         int32
    InetPrefix            string
}

type TotalConf struct {
    BasicC        BasicConf
    KitexConfig   KitexConf
	SaasOrgDB     MysqlConf
	ReqImgHost    ImgHost
}

func GetLocalHost() (string, error) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "", err
    }
    for _, a := range addrs {
        if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                myhost := ipnet.IP.String()
                return myhost, nil
            }
        }
    }
    return "", nil
}
