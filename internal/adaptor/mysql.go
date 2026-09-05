package adaptor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 生成 MySQL DSN, 用于连接 MySQL 数据库
func MySqlDSN(cfg config.MySQLConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)
}

// 建立连接池并套上 GORM：底层仍是 *sql.DB，返回 *gorm.DB 给仓储使用
func OpenMySql(cfg config.MySQLConfig) (*gorm.DB, error) {
	sqlDB, err := sql.Open("mysql", MySqlDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("打开 MySQL 失败: %w", err)
	}

	// 打开成功后，设置连接池大小
	maxOpen, maxIdle := cfg.MaxOpenConns, cfg.MaxIdleConns
	if maxOpen <= 0 {
		maxOpen = 10
	}

	if maxIdle <= 0 {
		maxIdle = 5
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)

	// 检查连接是否正常
	if err := PingMySQL(context.Background(), sqlDB); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	// 用现成的连接池套上 GORM，不再开第二条连接
	isMySQL := mysql.New(mysql.Config{Conn: sqlDB})
	gdb, err := gorm.Open(isMySQL, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印 SQL 语句
	})
	if err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("打开 GORM 失败: %w", err)
	}

	return gdb, nil
}

// 检查 MySQL 连接是否正常, 返回 error
func PingMySQL(ctx context.Context, db *sql.DB) error {
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping MySQL 失败: %w", err)
	}

	return nil
}
