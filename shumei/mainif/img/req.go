package img

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/tidwall/gjson"
	"img-tools/service/org_sign"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"sync/atomic"
)

type Img struct {

}

var ReqData struct {
	requestId string
	organization string
	url string
	callbackUrl string
	status       string
	reqParams   string
	retParams   string
}

func (this *Img) Predict(reqesutId string, reqParam ReqData) {
	start := time.Now().UnixNano()
	client := &http.Client{Timeout: 20 * time.Second}
	url := "http://" + t.Host + ":" + t.Port + request.Uri
	req, errIgnore := http.NewRequest("POST", url, bytes.NewBuffer([]byte(request.Data)))
	if errIgnore != nil {
		err = errIgnore
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("m_request_id", request.RequestId)
	req.Header.Set("m_organization", request.Organization)
	resp, errIgnore := client.Do(req)
	if errIgnore != nil {
		err = errIgnore
		return
	}
	defer resp.Body.Close()
	respBytes, errIgnore := ioutil.ReadAll(resp.Body)
	if errIgnore != nil {
		err = errIgnore
	}
	ret = string(respBytes)

	DiffResult(ret, reqParam.retParams)
	cost = utils.GetCost(start, time.Now().UnixNano())
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