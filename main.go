package main

import (
	"fmt"
	fileUtils "github.com/Hind3ight/FindProxy/pkg/file"
	"github.com/Hind3ight/FindProxy/pkg/ip"
	"strings"
)

func main() {
	content, err := fileUtils.OpenFile("./assets/originSource.txt")
	if err != nil {
		fmt.Printf("open file fail,err: %s\n", err)
	}
	cidrSlice := strings.Split(string(content), "\r\n")
	count := 0
	for _, cidr := range cidrSlice {
		ipPool := ip.ParseCidr(cidr)
		count += len(ipPool)
		//for _, v := range ipPool {
		//	fmt.Println(v)
		//	count++
		//}
	}
	fmt.Printf("总计%v条数据\n", count)
}
