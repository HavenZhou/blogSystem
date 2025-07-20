package database

import (
	"blogSystem/internal/domain"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(dsn string) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 		a) PrepareStmt: true
		// 启用预处理语句（Prepared Statement）缓存
		// 作用：
		// 提高性能：SQL 语句会被预编译并缓存，重复执行时直接使用编译好的语句
		// 防止 SQL 注入：自动对参数进行转义处理
		// 适用场景：高频执行的 SQL 语句

		// b) SkipDefaultTransaction: true
		// 禁用 GORM 的默认事务
		// 默认行为（false 时）：每个独立的写操作（Create/Update/Delete）都会自动包装在事务中
		// 设置为 true 后：写操作不再自动使用事务，需要手动管理事务提高性能（减少事务开销）;适合简单操作或需要自己控制事务的场景
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return err
	}

	sqlDB, err := DB.DB() // 用于获取这个底层连接池对象

	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)                  // 设置连接池中最大空闲连接数（默认值通常为 2）
	sqlDB.SetMaxOpenConns(100)                 // 设置最大打开连接数（默认无限制）
	sqlDB.SetConnMaxLifetime(time.Hour)        // 设置连接的最大存活时间（默认无限制）
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 设置连接最大空闲时间

	// 自动迁移
	if err := DB.AutoMigrate(
		&domain.User{}, &domain.Post{}, &domain.Comment{},
	); err != nil {
		return err
	}
	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close() // 关闭连接池
}

func GetDB() *gorm.DB {
	return DB
}
