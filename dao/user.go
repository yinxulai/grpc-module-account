package dao

import (
	"strconv"
	"strings"
	"time"

	"github.com/grpcbrick/account/models"
	"github.com/yinxulai/goutils/config"
	"github.com/yinxulai/goutils/crypto"
	"github.com/yinxulai/goutils/easysql"
)

const userTableName = "users"

func truncateUserTable() error {
	conn := easysql.GetConn()

	_, err := conn.ExecSQL("truncate table " + userTableName)
	return err
}

func createUserTable() error {
	conn := easysql.GetConn()

	_, err := conn.ExecSQL(
		strings.Join([]string{
			" CREATE TABLE IF NOT EXISTS `" + userTableName + "`(",
			" `ID` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',",
			" `Class` varchar(128) NOT NULL COMMENT '账户类型',",
			" `Avatar` varchar(512) DEFAULT '' COMMENT '头像', ",
			" `Inviter` int(11) DEFAULT 0 COMMENT '邀请人',",
			" `Nickname` varchar(128) NOT NULL COMMENT '昵称',",
			" `Username` varchar(128) NOT NULL COMMENT '用户名',",
			" `Password` varchar(512) NOT NULL COMMENT '密码',",
			" `DeletedTime` datetime DEFAULT NULL COMMENT '删除时间',",
			" `CreatedTime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',",
			" `UpdatedTime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',",
			" PRIMARY KEY (`ID`,`Nickname`,`Class`,`Username`),",
			" UNIQUE KEY `Username` (`Username`)",
			" ) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;",
		}, "",
		),
	)
	return err
}

// CountUserByID 根据 id 统计
func CountUserByID(id int64) (int, error) {
	conn := easysql.GetConn()

	idstr := strconv.FormatUint(uint64(id), 10)
	cond := map[string]string{"ID": idstr}
	queryField := []string{"count(*) as count"}
	result, err := conn.Select(userTableName, queryField).Where(cond).QueryRow()
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(result["count"])
	if err != nil {
		return 0, err
	}
	return count, nil
}

// QueryUsers 查询用户
func QueryUsers(page, limit int) (totalPage, currentPage int, users []*models.User, err error) {
	// conn := easysql.GetConn()
	// totalPage, currentPage, rows, err := conn.Select(userTableName, nil).Pagination(page, limit)
	// if err != nil {
	// 	return 0, 0, nil, err
	// }

	// users = []*models.User{}
	// for _, mapData := range rows {
	// 	user := new(models.User)
	// 	user.LoadStringMap(mapData.(map[string]string))
	// 	users = append(users, user)
	// }

	return totalPage, currentPage, users, err
}

// QueryUsersByInviter 查询用户
func QueryUsersByInviter(inviter int64, page, limit int) (totalPage, currentPage int, users []*models.User, err error) {
	// conn := easysql.GetConn()
	// inviterstr := strconv.FormatUint(uint64(inviter), 10)
	// cond := map[string]string{"Inviter": inviterstr}
	// totalPage, currentPage, result, err := conn.Select(userTableName, nil).Where(cond).Pagination(page, limit)
	// if err != nil {
	// 	return 0, 0, nil, err
	// }

	// users = []*models.User{}
	// for _, mapData := range result {
	// 	user := new(models.User)
	// 	user.LoadStringMap(mapData.(map[string]string))
	// 	users = append(users, user)
	// }
	return totalPage, currentPage, users, nil
}

// QueryUserByID 根据 id 查询
func QueryUserByID(id int64) (*models.User, error) {
	conn := easysql.GetConn()

	idstr := strconv.FormatUint(uint64(id), 10)
	cond := map[string]string{"ID": idstr}
	result, err := conn.Select(userTableName, nil).Where(cond).QueryRow()
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	user.LoadStringMap(result)
	return user, nil
}

// QueryUserByUsername 根据 id 查询
func QueryUserByUsername(username string) (*models.User, error) {
	conn := easysql.GetConn()

	cond := map[string]string{"Username": "'" + username + "'"}
	result, err := conn.Select(userTableName, nil).Where(cond).QueryRow()
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	user.LoadStringMap(result)
	return user, nil
}

// CountUserByUsername 根据用户名统计
func CountUserByUsername(username string) (int, error) {
	conn := easysql.GetConn()

	queryField := []string{"count(*) as count"}
	cond := map[string]string{"Username": "'" + username + "'"}
	result, err := conn.Select(userTableName, queryField).Where(cond).QueryRow()
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(result["count"])
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateUser 创建用户
func CreateUser(class, nickname, username, password string, inviter int64) (int64, error) {
	conn := easysql.GetConn()

	cond := map[string]string{
		"Class":    class,
		"Nickname": nickname,
		"Username": username,
		"Inviter":  strconv.FormatUint(uint64(inviter), 10),
		"Password": crypto.MD5Encrypt(password, config.MustGet("encrypt-password")),
	}

	id, err := conn.Insert("users", cond)
	if err != nil {
		return id, err
	}
	return id, nil
}

// DeleteUserByID 删除用户
func DeleteUserByID(id int64) error {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	return UpdataUserFieldByID(id, map[string]string{"DeletedTime": nowTime})
}

// UpdataUserFieldByID 根据 ID 更新用户指定字段
func UpdataUserFieldByID(id int64, field map[string]string) error {
	conn := easysql.GetConn()

	cond := map[string]string{"ID": strconv.FormatUint(uint64(id), 10)}
	_, err := conn.Where(cond).Update(userTableName, field)
	return err
}

// UpdateUserClassByID 更新用户类型
func UpdateUserClassByID(id int64, class string) error {
	return UpdataUserFieldByID(id, map[string]string{"Class": class})
}

// UpdateUserAvatarByID 更新用户头像
func UpdateUserAvatarByID(id int64, avatar string) error {
	return UpdataUserFieldByID(id, map[string]string{"Avatar": avatar})
}

// UpdateUserNicknameByID 更新用户昵称
func UpdateUserNicknameByID(id int64, nickname string) error {
	return UpdataUserFieldByID(id, map[string]string{"Nickname": nickname})
}

// UpdateUserInviterByID 更新用户邀请码
func UpdateUserInviterByID(id, inviter int64) error {
	return UpdataUserFieldByID(id, map[string]string{"Inviter": strconv.FormatUint(uint64(inviter), 10)})
}

// UpdateUserPasswordByID 更新用户密码
func UpdateUserPasswordByID(id int64, password string) error {
	// 加密
	encryptPassword := crypto.MD5Encrypt(password, config.MustGet("encrypt-password"))
	return UpdataUserFieldByID(id, map[string]string{"Password": encryptPassword})
}

// VerifyUserPasswordByID 验证用户密码
func VerifyUserPasswordByID(id int64, password string) (bool, error) {
	conn := easysql.GetConn()

	idstr := strconv.FormatUint(uint64(id), 10)
	cond := map[string]string{"ID": idstr}
	result, err := conn.Select(userTableName, nil).Where(cond).QueryRow()
	if err != nil {
		return false, err
	}

	// 加密
	encryptPassword := crypto.MD5Encrypt(password, config.MustGet("encrypt-password"))
	if result["Password"] == encryptPassword {
		return true, nil
	}

	return false, nil
}
