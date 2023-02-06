package dao

import (
	"database/sql"
)

// QueryByUserName 从user表中根据userName字段条件查询
func QueryByUserName(userName string) *sql.Row {
	sqlStr := "select * from user where userName=?;"
	row := db.QueryRow(sqlStr, userName)

	return row
}

// InsertUser 插入一条记录到user表
func InsertUser(userName string, password string) error {
	sql := "insert into user values(?,?)"
	_, err := db.Exec(sql, userName, password)

	return err
}

// QueryUser 从user表中根据userName, password字段条件查询
func QueryUser(userName string, password string) *sql.Row {
	sql := "select * from user where userName=? and password=?"
	row := db.QueryRow(sql, userName, password)

	return row
}

// ExecInsert 执行sql插入语句
func ExecInsert(sql string) (sql.Result, error) {
	result, err := db.Exec(sql)

	return result, err
}

// ExecUpdate 执行sql更新语句
func ExecUpdate(sql string) (sql.Result, error) {
	return ExecInsert(sql)
}

// ExecQueryRow 执行sql查询语句，返回一个结果
func ExecQueryRow(sql string) *sql.Row {
	rows := db.QueryRow(sql)

	return rows
}

// ExecQueryRows 执行sql查询语句，返回结果集
func ExecQueryRows(sql string) (*sql.Rows, error) {
	rows, err := db.Query(sql)

	return rows, err
}
