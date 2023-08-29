package img

import (
	"bytes"
	//"encoding/json"
	//"errors"
	"code.aliyun.com/module-go/protocols/image/kitex_gen/shumei/strategy/re"
	"github.com/google/go-cmp/cmp"
    "shumei/mainif/config"
    "shumei/mainif/rpc"
	//"github.com/tidwall/gjson"
	"io/ioutil"
	"math"
    "fmt"
	"net/http"
	"strings"
    "time"
	//"sync/atomic"
)

type Img struct {

}

type ReqData struct {
	RequestId    string
	Organization string
	Url          string
	CallbackUrl  string
	Status       string
	ReqParams    string
	RetParams    string
}

func (this *Img) Predict(reqesutId string, reqParam ReqData) {
	start := time.Now().UnixNano()
	reqParam := re.ImagePredictRequest{}
	respBytes, errr := rpc.reImageClientModel.Predict(&reqParam)
	rpc.reImageClientModel.Version()
	ret := string(respBytes)
    fmt.Println("rrrrrrrrr:", ret)
	this.DiffResult(ret, reqParam.RetParams)
	cost := this.GetCost(start, time.Now().UnixNano())
    fmt.Println("ret1 and ret2 and request cost:", ret, reqParam.RetParams, cost)
	return
}


func (this *Img) DiffResult(data1 string, data2 string) string {
	if strings.Contains(data1, "参数不合法") {
		return ""
	}

	if strings.Contains(data2, "参数不合法") {
		return ""
	}
	return cmp.Diff(data1, data2, cmp.Comparer(func(a, b float64) bool {
		switch {
		case math.Abs(a-b) <= 0.001:
			return true
		}
		return false
	}))
}

func (this *Img) GetCost(startTime int64, endTime int64) float64 {
    return float64(endTime-startTime) / float64(1000000)
}

