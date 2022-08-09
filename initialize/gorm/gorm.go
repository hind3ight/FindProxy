//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package gorm

import (
	"fmt"
	"github.com/Hind3ight/FindProxy/global"
	"github.com/Hind3ight/FindProxy/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	userName := "root"
	password := "123456"
	path := "127.0.0.1:3306"
	dbName := "proxyPool"
	logMode := logger.Info
	// root:123456@tcp( 127.0.0.1:3306)/proxyPool?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", userName, password, path, dbName)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		panic(err)
	} else {
		Migrate(db)
	}

	global.DB = db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(model.AllIp{}, model.OriginIp{})
	if err != nil {
		fmt.Println("自动迁移模型失败")
		panic(err)
	}
}
