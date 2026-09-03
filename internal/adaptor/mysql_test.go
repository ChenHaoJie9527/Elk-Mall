package adaptor

import (
	"context"
	"testing"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/config"
)

func TestMySqlDSN(t *testing.T) {
	cfg := config.MySQLConfig{
		Host:     "127.0.0.1",
		Port:     "3308",
		User:     "root",
		Password: "root",
		Database: "elk_mall",
	}
	want := "root:root@tcp(127.0.0.1:3308)/elk_mall?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := MySqlDSN(cfg)
	if dsn != want {
		t.Fatalf("MySQL DSN 不正确: 期望 %v, 实际 %v", want, dsn)
	}
	t.Logf("MySQL DSN 测试通过: %v", dsn)
}

func TestOpenMySql(t *testing.T) {
	cfg := config.MySQLConfig{
		Host:     "127.0.0.1",
		Port:     "3308",
		User:     "root",
		Password: "root",
		Database: "elk_mall",
	}
	db, err := OpenMySql(cfg)
	if err != nil {
		t.Fatalf("打开 MySQL 失败: %v", err)
	}
	defer db.Close()

	if err := PingMySQL(context.Background(), db); err != nil {
		t.Fatalf("ping MySQL 失败: %v", err)
	}

	t.Logf("MySQL 连接成功: %v", db.Stats())
}
