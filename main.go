package main

import (
	"fmt"
	_ "github.com/Hind3ight/FindProxy/initialize/gorm"
	"github.com/Hind3ight/FindProxy/model"
	fileUtils "github.com/Hind3ight/FindProxy/pkg/file"
	"github.com/Hind3ight/FindProxy/pkg/ip"
	"github.com/Hind3ight/FindProxy/service"
	"strings"
	"time"
)

var dbDataList = make(chan interface{}, 1000)

func main() {
	go MonitorAndSave()

	content, err := fileUtils.OpenFile("./assets/originSource.txt")
	if err != nil {
		fmt.Printf("open file fail,err: %s\n", err)
	}
	cidrSlice := strings.Split(string(content), "\r\n")
	var count1, count2 = 0, 0
	var dbData []model.OriginIp
	for _, cidr := range cidrSlice {
		ipPool := ip.ParseCidr(cidr)
		count1 += len(ipPool)

		for _, v := range ipPool {
			one := model.OriginIp{
				Ip: v,
			}
			dbData = append(dbData, one)
			count2++
			if count2 == 1000 {
				dbDataList <- dbData
				count2 = 0
				dbData = nil
			}
		}
	}
	fmt.Printf("总计%v条数据\n", count1)
}

func MonitorAndSave() {
	ipSrv := service.IpSrv{}

	for {
		select {
		case dbData := <-dbDataList:
			ipSrv.SaveIp(dbData.([]model.OriginIp))
		default:
			time.Sleep(time.Second * 5)

		}
	}
}
