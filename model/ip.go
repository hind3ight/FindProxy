//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package model

import (
	"gorm.io/gorm"
	"time"
)

type AllIp struct {
	ID            uint64    `json:"id" form:"id" gorm:"column:id;comment:主键"`
	CreatedAt     time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;comment:创建时间"`
	UpdateAt      time.Time `json:"updateAt" form:"updateAt" gorm:"column:update_at;comment:更新时间"`
	DeleteAt      time.Time `json:"deleteAt" form:"deleteAt" gorm:"column:delete_at;comment:删除时间"`
	Ip            string    `json:"ip" form:"ip" gorm:"column:ip;comment:"`
	LastCheckDate time.Time `json:"lastCheckDate" form:"lastCheckDate" gorm:"column:last_check_date;comment:上次检测时间"`
	IsEffective   bool      `json:"isEffective" form:"isEffective" gorm:"column:is_effective;comment:是否有效"`
}

type OriginIp struct {
	gorm.Model
	Ip string `json:"ip" form:"ip" gorm:"column:ip;comment:"`
}
