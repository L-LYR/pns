/*
	test data generator
*/

package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/L-LYR/pns/mobile/push_sdk/net/http"
)

const (
	DispatchN = 8
)

var (
	bizClient          = http.MustNewHTTPClient("http://127.0.0.1:10087")
	client             = http.MustNewHTTPClient("http://127.0.0.1:10086")
	osList             = []string{"windows", "android", "macos", "linux"}
	brandList          = []string{"chrome", "huawei", "vivo", "apple", "safari", "firefox", "windows"}
	modelList          = []string{"xxxx", "yyyy", "zzzz"}
	appVersionList     = []string{"0.0.1", "0.2.1", "0.2.2", "0.1.0", "0.1.1"}
	pushSDKVersionList = []string{"0.0.1", "0.0.2", "0.0.3"}
)

func getRandTarget(appId int, deviceId int64) http.Payload {
	return http.Payload{
		"deviceId":           deviceId,
		"os":                 osList[rand.Intn(len(osList))],
		"brand":              brandList[rand.Intn(len(brandList))],
		"model":              modelList[rand.Intn(len(modelList))],
		"tzName":             "Asia/Shanghai",
		"appId":              appId,
		"appVersion":         appVersionList[rand.Intn(len(appVersionList))],
		"pushSDKVersion":     pushSDKVersionList[rand.Intn(len(pushSDKVersionList))],
		"language":           "cn",
		"inAppPushStatus":    rand.Intn(2),
		"systemPushStatus":   rand.Intn(2),
		"privacyPushStatus":  rand.Intn(2),
		"businessPushStatus": make(map[string]int),
	}
}

func mockUpdateTarget(appId int, deviceId int64) error {
	_, err := client.POST("/target", getRandTarget(appId, deviceId))
	return err
}

func mockCreateApp(appId int) error {
	_, err := bizClient.POST("/app", map[string]interface{}{
		"appId":   appId,
		"appName": fmt.Sprintf("testapp%d", appId),
	})
	return err
}

func mockOpenMQTT(appId int) error {
	_, err := bizClient.POST("/app/config/mqtt", map[string]interface{}{
		"appId": appId,
	})
	return err
}

func GenDataSet(appId int, maxDeviceNum int64) {
	if err := mockCreateApp(appId); err != nil {
		log.Printf("app %d, err=%+v", appId, err)
		return
	}

	if err := mockOpenMQTT(appId); err != nil {
		log.Printf("app %d, err=%+v", appId, err)
		return
	}

	var wg sync.WaitGroup

	counter := int64(0)
	reqNumPerWorker := maxDeviceNum / DispatchN

	wg.Add(DispatchN)
	for i := 0; i < DispatchN; i++ {
		go func(begin int64, id int) {
			log.Printf("%d worker\n", id)
			for i := begin; i <= begin+reqNumPerWorker; i++ {
				if err := mockUpdateTarget(appId, i); err != nil {
					log.Printf("device %d, err=%+v", i, err)
				}
				atomic.AddInt64(&counter, 1)
			}
			wg.Done()
		}(int64(i)*reqNumPerWorker, i)
	}

	wg.Add(1)
	go func() {
		for range time.Tick(time.Second) {
			c := atomic.LoadInt64(&counter)
			log.Printf("%d/%d\n", c, maxDeviceNum)
			if c >= maxDeviceNum {
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
}

func main() {
	rand.Seed(time.Now().Unix())
	/*
		AppID 		DeviceNumber
		1234  		  1000000
		12341 		  100000
		12342 		  100000
		12343 		  100000
		12344 		  100000
		12345 		  100000
		12346 		  100000
		12347 		  100000
		12348 		  100000
		12349 		  100000
		123411		  10000
		 ...
		123420		  10000
		123421		  5000
		 ...
		123450		  5000
		123451        1000
		 ...
		123490        1000
	*/
	GenDataSet(1234, 1000000)
	for i := 12341; i <= 12349; i++ {
		GenDataSet(i, 100000)
	}
	for i := 123411; i <= 123420; i++ {
		GenDataSet(i, 10000)
	}
	for i := 123421; i <= 123450; i++ {
		GenDataSet(i, 5000)
	}
	for i := 123451; i <= 123490; i++ {
		GenDataSet(i, 1000)
	}
}
