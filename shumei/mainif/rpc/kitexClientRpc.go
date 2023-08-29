package rpc

import (
	"code.aliyun.com/module-go/prometheus"
	"code.aliyun.com/module-go/protocols/image/kitex_gen/shumei/strategy/re"
	"code.aliyun.com/module-go/protocols/image/kitex_gen/shumei/strategy/re/imagepredictor"
	"code.aliyun.com/module-go/registry-zookeeper/resolver"
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"time"
)

const maxRpcTotalTimeOut = 7000

// reImageClientModel reImageClient 单例对象，持有reImageClient
type reImageClientModel struct {
	ReImageClient *imagepredictor.Client
	//RpcConfig     *rpcModel.RpcConfigParam
}

// Predict reImageClient真实请求，对外暴露对接口
func (clientMode *reImageClientModel) Predict(predictRequest *re.ImagePredictRequest) (*re.ImagePredictResult_, error) {
	reImageClient := *clientMode.ReImageClient
	if reImageClient != nil {
		predictResult, err := reImageClient.Predict(context.Background(), predictRequest)
		return predictResult, err
	}
	return nil, errors.New("reImage client is nil")
}

func (clientMode *reImageClientModel) Version() (resp string, err error) {
	version, err := reImageClient.Version(context.Background())
	fmt.Println("invoke version func:", version, err)
	return version, err
}

// init 初始化reImageClient 对象，提供给工厂init，参数校验，初始化
func (clientMode *reImageClientModel) init() {
	clientMode.initReImageClient()
	return 
}

//初始化相关逻辑,真正初始化逻辑
func (clientMode *reImageClientModel) initReImageClient()  {
	start := time.Now().UnixNano()
	utils := common.Utils{}
	r, err := resolver.NewZookeeperResolver(clientMode.getZkAddressList(), clientMode.getZkSessionTime()*time.Millisecond)
	if err != nil {
		fmt.Println("NewZookeeperResolver err:", err)
		return 
	}
	policy := retry.NewFailurePolicy()
	policy.ShouldResultRetry = retry.AllErrorRetry()
	policy.WithMaxRetryTimes(clientMode.getMaxRetryTimes())
	policy.WithRetryBreaker(clientMode.getRetryBreaker())
	policy.WithMaxDurationMS(maxRpcTotalTimeOut)
	reImageClient, err := imagepredictor.NewClient(clientMode.getRpcServicePath(), client.WithResolver(r),
		client.WithRPCTimeout(clientMode.getRpcTimeout()*time.Millisecond),
		client.WithShortConnection(),
		client.WithTracer(prometheus.NewClientTracerWithoutExport()),
		client.WithConnectTimeout(clientMode.getRpcConnectTimeout()*time.Millisecond),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "re-image-ng"}),
		client.WithFailureRetry(policy))
	if err != nil {
		return 
	}
	version, err := reImageClient.Version(context.Background())
	fmt.Println("version:", version, ",err:", err)
	if err != nil {
		fmt.Println("version err:", version, ",err:", err)
	}

	clientMode.setReImageClient(&reImageClient)
	return 
}


func (clientMode reImageClientModel) getRpcServicePath() string {
	servicePath := "/request-engine/re-image-ng"//clientMode.RpcConfig.ServicePath
	return servicePath
}

//获取zkAddress
func (clientMode reImageClientModel) getZkAddressList() []string {
	temp := [1]{"10.141.4.128:2181"}
	return temp//clientMode.RpcConfig.ZkAddressList
}

//设置reImageClient
func (clientMode *reImageClientModel) setReImageClient(reImageClient *imagepredictor.Client) {
	clientMode.ReImageClient = reImageClient
}

//获取reImageClient
func (clientMode reImageClientModel) getReImageClient() *imagepredictor.Client {
	return clientMode.ReImageClient
}

//获取zkSessionTime
func (clientMode reImageClientModel) getZkSessionTime() time.Duration {
	return 6000//clientMode.RpcConfig.ZkSessionTime
}

//获取rpcConnectionOUt
func (clientMode reImageClientModel) getRpcConnectTimeout() time.Duration {
	return 6000//clientMode.RpcConfig.RpcConnectTimeout
}

//获取最大重试次数
func (clientMode reImageClientModel) getMaxRetryTimes() int {
	return 1//clientMode.RpcConfig.MaxRetryTimes
}

//获取熔断率
func (clientMode reImageClientModel) getRetryBreaker() float64 {
	return 0.2//clientMode.RpcConfig.RetryBreaker
}

//获取rpc请求超时时间
func (clientMode reImageClientModel) getRpcTimeout() time.Duration {
	return 3000//clientMode.RpcConfig.RpcTimeout
}
