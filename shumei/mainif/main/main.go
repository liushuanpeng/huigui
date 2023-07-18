package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
)

var Db *sql.DB
var configFilePath *string

func commandLine() {
	configFilePath = flag.String("cfg", "./config.json", "配置文件路径，默认使用==> ./config.json")
	flag.Parse()
}

func initDB(myc *config.MysqlConf) (*sql.DB, error) {
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%dms&readTimeout=%dms&writeTimeout=%dms&charset=utf8", myc.User, myc.Password, myc.Host, myc.Port, myc.DbName, myc.ConnTimeout, myc.ReadTimeout, myc.WriteTimeout)
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("Mysql Connection error")
	}
	db.SetMaxOpenConns(myc.MaxOpenConn)
	db.SetMaxIdleConns(myc.MaxIdleConn)
	db.Ping()
	return db, nil
}

func initDBs() error {

	var db *sql.DB
	var err error
	if db, err = initDB(config.Conf.ConfigMap.SaasOrgDB); err != nil {
		return err
	}
	Db = db
}

func loadConfig() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error:", err)
			debug.PrintStack()
		}
	}()
	config.Conf = new(config.Config)
	err := config.Conf.LoadConfig(*configFilePath)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	commandLine()
	err := loadConfig()
	if err != nil {
		fmt.Println("loadConfig error", err)
		return 
	}
	initDBs()

	//读取数据库
	retParam := ReadMysql()
	//发起请求
	img := img.Img{}
	for key, value := range img {
	    img.Predict(requestId, value)
	}
	return
}


func  ReadMysql() map[string]interface{} {
	excu := fmt.Sprintf("select `requestId`, `org`, `url`, `callback_url`, `status`, `req_params`, `ret_params` from pi_req_and_ret_params where `status` = %d", 1)
	rows, err := Db.Query(excu)
	if err != nil {
		fmt.Println("query err:", err)
		return nil
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("rows err:", err)
			return nil
		}
	}()

	var requestId string
	var organization string
	var url string
	var callback_url string
	var status       string
	var req_params   string
	var ret_params   string
	TempData := make(map[string]interface{})
	for rows.Next() {
		if rowScanErr := rows.Scan(&requestId, &organization, &url, &callback_url, &status, &req_params, &ret_params); rowScanErr != nil {
			fmt.Println("rows scan err:", rowScanErr)
			continue
		}
		reqData := img.ReqData{requestId: requestId, organization:organization, url:url, callbackUrl:callback_url, status: status, reqParams: req_params, retParams:ret_params}
		TempData[requestId] = reqData
	}
	return TempData
}