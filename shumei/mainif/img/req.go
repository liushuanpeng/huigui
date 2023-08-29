package img

import (
	"bytes"
	//"encoding/json"
	//"errors"
	"github.com/google/go-cmp/cmp"
    "shumei/mainif/config"
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
	client := &http.Client{Timeout: 20 * time.Second}
	url := "http://" + config.Conf.ConfigMap.ReqImgHost.Host + ":" + config.Conf.ConfigMap.ReqImgHost.Port + reqParam.Url
	req, errIgnore := http.NewRequest("POST", url, bytes.NewBuffer([]byte(reqParam.ReqParams)))
	if errIgnore != nil {
		fmt.Println("NewRequest errIgnore", errIgnore)
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("m_request_id", reqParam.RequestId)
	req.Header.Set("m_organization", reqParam.Organization)
	resp, errIgnore := client.Do(req)
	if errIgnore != nil {
		fmt.Println("client.Do Err:", errIgnore)
		return
	}
	defer resp.Body.Close()
	respBytes, errIgnore := ioutil.ReadAll(resp.Body)
	if errIgnore != nil {
        fmt.Println("ReadAll Err:", errIgnore)
        return
	}
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

