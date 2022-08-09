//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package service

import (
	"fmt"
	"github.com/Hind3ight/FindProxy/global"
	"github.com/Hind3ight/FindProxy/model"
)

type IpSrv struct {
}

func (s IpSrv) SaveIp(dbData []model.OriginIp) {
	err := global.DB.Save(&dbData).Error
	if err != nil {
		fmt.Println(err)
	}
}
