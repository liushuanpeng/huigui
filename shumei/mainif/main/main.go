package main

import (
    _ "github.com/go-sql-driver/mysql"
    "shumei/mainif/config"
    "shumei/mainif/img"
    "shumei/mainif/rpcClient"
	"database/sql"
    "runtime/debug"
    "strings"
    "time"
    "net"
	"errors"
	"flag"
	"fmt"
)

var Db *sql.DB
var configFilePath *string
var port *int64
var r registry.Registry
var info *registry.Info

func commandLine() {
	configFilePath = flag.String("cfg", "./config.json", "配置文件路径，默认使用==> ./config.json")
    port           = flag.Int64("port", 7888, "port")
	flag.Parse()
}

func initDB(myc config.MysqlConf) (*sql.DB, error) {
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
    return nil
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

	reImageClientMode := reImageClientModel{}
	reImageClientMode.init()

    go rpc.kitexRpcServer()
    time.Sleep(time.Second*5)
	err = initDBs()
    if err != nil {
        fmt.Println("initDBs err:", err)
        return    
    }

	//读取数据库
	retParam := ReadMysql()
	//发起请求
	img := img.Img{}
	for key, value := range retParam {
	    img.Predict(key, value)
	}
	return
}

func  ReadMysql() map[string]img.ReqData {
	excu := fmt.Sprintf("select `requestId`, `organization`, `url`, `callback_url`, `status`, `req_params`, `ret_params` from pi_req_and_ret_params where `status` = %d", 1)
	rows, err := Db.Query(excu)
	if err != nil {
		fmt.Println("query err:", err)
		return nil
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("rows err:", err)
			return
		}
	}()

	var requestId string
	var organization string
	var url string
	var callback_url string
	var status       string
	var req_params   string
	var ret_params   string
	TempData := make(map[string]img.ReqData)
	for rows.Next() {
		if rowScanErr := rows.Scan(&requestId, &organization, &url, &callback_url, &status, &req_params, &ret_params); rowScanErr != nil {
			fmt.Println("rows scan err:", rowScanErr)
			continue
		}
		reqData := img.ReqData{RequestId: requestId, Organization:organization, Url:url, CallbackUrl:callback_url, Status: status, ReqParams: req_params, RetParams:ret_params}
		TempData[requestId] = reqData
	}
	return TempData
}
