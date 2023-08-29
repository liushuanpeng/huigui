package rpc

import (
    zkRegistry "code-int.ishumei.com/module-go/registry-zookeeper/registry"
    re "code-int.ishumei.com/module-go/protocols/image/kitex_gen/shumei/strategy/re/imagepredictor"
    "shumei/mainif/config"
}

func kitexRpcServer() {
    host, errH := config.GetLocalHost()
    if errH != nil || host == "" {
        fmt.Println("err get local host:", errH, host)
        return
    }

    addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%d", host, *port))
    zkPath := "request-engine/re-image-ng"//config.Conf.ConfigMap.BasicC.ZkPath[1:] //去掉前面的/
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