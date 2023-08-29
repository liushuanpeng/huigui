package main

import (
    zkRegistry "github.com/liushuanpeng/huigui/registry-zookeeper/registry"
    re "github.com/liushuanpeng/huigui/protocols/image/kitex_gen/shumei/strategy/re/imagepredictor"
    _ "github.com/go-sql-driver/mysql"
    "shumei/mainif/config"
    "shumei/mainif/img"
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
    go kitexRpcServer()
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


func kitexRpcServer() {
    host, errH := config.GetLocalHost()
    if errH != nil || host == "" {
        fmt.Println("err get local host:", errH, host)
        return
    }

    addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%d", host, *port))
    zkPath := config.Conf.ConfigMap.BasicC.ZkPath[1:] //去掉前面的/
    //zk server
    var err error
    zkServers := strings.Split(config.Conf.ConfigMap.BasicC.ZkServers, ",")
    r, err = zkRegistry.NewZookeeperRegistry(zkServers, 40*time.Second)
    if err != nil {
        fmt.Println("register err:", err)
    }
    //tags := map[string]string{"group": "blue", "idc": "hd1"}
    weight := config.Conf.ConfigMap.KitexConfig.Weight
    maxCon := config.Conf.ConfigMap.KitexConfig.MaxConnections
    maxQps := config.Conf.ConfigMap.KitexConfig.MaxQPS
    if weight <= 0 {
        weight = 100
    }
    if maxCon <= 0 {
        maxCon = 10000
    }
    if maxQps <= 0 {
        maxQps = 2000
    }
    fmt.Println("kitex config weight, maxCon, maxQps:", weight, maxCon, maxQps)
    info = &registry.Info{ServiceName: zkPath, Weight: weight, PayloadCodec: "thrift", Addr: addr}

    svr := re.NewServer(
        new(predict.ImagePredictorImpl),
        server.WithServiceAddr(addr), //&net.TCPAddr{IP:net.ParseIP(host),Port:int(*port)}),
        server.WithRegistry(r),
        //server.WithRegistryInfo(info),
        //server.WithRegistry(models.DefaultRegistry),
        //server.WithTracer(prometheus.NewServerTracer(fmt.Sprintf(":%d", *port+1), "/metrics")),
        //server.WithSuite(tracing.NewServerSuite()),
        server.WithLimit(&limit.Option{MaxConnections: maxCon, MaxQPS: maxQps}),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: zkPath}),
        server.WithTracer(prometheus.NewServerTracerWithoutExport()),
    )

    err = svr.Run()
    if err != nil {
        fmt.Println("run server err:", err.Error())
    }
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
