package adaptor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

// 生成 MySQL DSN, 用于连接 MySQL 数据库
func MySqlDSN(cfg config.MySQLConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)
}

// 建立连接池，设置池子大小，以及 ping 连接
func OpenMySql(cfg config.MySQLConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", MySqlDSN(cfg))
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

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)

	// 检查连接是否正常
	if err := PingMySQL(context.Background(), db); err != nil {
		// 如果连接不正常，关闭连接，并返回错误
		_ = db.Close()
		return nil, err
	}

	return db, nil

}

// 检查 MySQL 连接是否正常, 返回 error
func PingMySQL(ctx context.Context, db *sql.DB) error {
	// 使用 context 检查连接是否正常
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping MySQL 失败: %w", err)
	}

	return nil
}
