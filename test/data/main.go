/*
	test data generator
*/

package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/L-LYR/pns/mobile/push_sdk/net/http"
)

const (
	MaxDeviceNum    = 5000000
	DispatchN       = 50
	ReqNumPerWorker = MaxDeviceNum / DispatchN
)

var (
	client             = http.MustNewHTTPClient("http://192.168.1.2")
	osList             = []string{"windows", "android", "macos", "linux"}
	brandList          = []string{"chrome", "huawei", "vivo", "apple", "safari", "firefox", "windows"}
	modelList          = []string{"xxxx", "yyyy", "zzzz"}
	appVersionList     = []string{"0.0.1", "0.2.1", "0.2.2", "0.1.0", "0.1.1"}
	pushSDKVersionList = []string{"0.0.1", "0.0.2", "0.0.3"}
)

func getRandTarget(deviceId int) http.Payload {
	return http.Payload{
		"deviceId":           deviceId,
		"os":                 osList[rand.Intn(len(osList))],
		"brand":              brandList[rand.Intn(len(brandList))],
		"model":              modelList[rand.Intn(len(modelList))],
		"tzName":             "Asia/Shanghai",
		"appId":              12345,
		"appVersion":         appVersionList[rand.Intn(len(appVersionList))],
		"pushSDKVersion":     pushSDKVersionList[rand.Intn(len(pushSDKVersionList))],
		"language":           "cn",
		"inAppPushStatus":    1,
		"systemPushStatus":   1,
		"privacyPushStatus":  1,
		"businessPushStatus": make(map[string]int),
	}
}

func mockUpdateTarget(deviceId int) error {
	_, err := client.POST("/target", getRandTarget(deviceId))
	return err
}

func main() {
	rand.Seed(time.Now().Unix())

	var wg sync.WaitGroup
	wg.Add(DispatchN)
	for i := 0; i < DispatchN; i++ {
		begin := i * ReqNumPerWorker
		go func(begin int) {
			for i := begin; i <= begin+ReqNumPerWorker; i++ {
				if err := mockUpdateTarget(i); err != nil {
					log.Printf("device %d, err=%+v", i, err)
				}
			}
			wg.Done()
		}(begin)
	}
	wg.Wait()
}
