package dao

import (
	"testing"

	"github.com/yinxulai/goutils/easysql"
)

func TestCreateUserTable(t *testing.T) {
	easysql.Init("mysql", "root:root@tcp(localhost:3306)/default?charset=utf8")
	createUserTable()
	t.Error("test")
}

func TestCreateUser(t *testing.T) {
	easysql.Init("mysql", "root:root@tcp(localhost:3306)/default?charset=utf8")
	createUser("test", "test", "test", "test")
	t.Error("test")
}

func TestQueryUserByID(t *testing.T) {
	easysql.Init("mysql", "root:root@tcp(localhost:3306)/default?charset=utf8")
	queryUserByID(1)
	t.Error("test")
}

func TestQueryUserByUsername(t *testing.T) {
	easysql.Init("mysql", "root:root@tcp(localhost:3306)/default?charset=utf8")
	queryUserByUsername("")
	t.Error("test")
}
