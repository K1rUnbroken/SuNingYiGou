package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "golang.org/x/crypto/ssh"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"su-ning-yi-gou/dao"
	"su-ning-yi-gou/model"
	"time"
)

// IsRepeat 检查用户名是否重复
func IsRepeat(userName string) (error, bool) {
	row := dao.QueryByUserName(userName)
	var u model.User
	err := row.Scan(&u.UserName, &u.Password)

	// 没有查询到结果
	if err == sql.ErrNoRows {
		return nil, false
	}

	// 其他错误
	if err != nil {
		return err, true
	}

	return nil, true
}

// AddUser 添加用户到数据库
func AddUser(userName string, password string) error {
	err := dao.InsertUser(userName, password)

	return err
}

// IsExists 检查用户是否存在
func IsExists(userName string) (error, bool) {
	return IsRepeat(userName)
}

// IsRight 检查密码是否正确
func IsRight(userName string, password string) (error, bool) {
	row := dao.QueryUser(userName, password)
	var user model.User
	err := row.Scan(&user.UserName, &user.Password)

	// 没有查询到结果
	if err == sql.ErrNoRows {
		return nil, false
	}

	// 其他错误
	if err != nil {
		return err, true
	}

	return nil, true
}

// JWTExpireDuration JWT过期时间
const JWTExpireDuration = time.Hour * 1

// JWT密匙
var secret = []byte("red rock")

// GenJWT 生成JWT
func GenJWT(userName string) (string, error) {
	// 创建一个claims，相当于jwt的payload
	c := model.Claims{
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(JWTExpireDuration).Unix(),
			Issuer:    "kyr",
		},
	}

	// 创建一个jwt对象
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// 把jwt的header和payload进行编码后再加上我们的secret，使用header中指定的算法进行签名，返回字符串
	return jwt.SignedString(secret)
}

// ParseJWT 解析JWT
func ParseJWT(jwtStr string) (*model.Claims, error) {
	// 解析
	jwt, err := jwt.ParseWithClaims(jwtStr, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	// 校验
	claims, isOK := jwt.Claims.(*model.Claims)
	if !(isOK && jwt.Valid) {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetProdByNameAndType 根据 商品名和商品类型 模糊查询，返回指定页数的结果，不包括商品评论
func GetProdByNameAndType(prodKeyWords string, offset int, pageSize int, sort string) ([]model.Product, error) {
	// 创建sql语句
	var sql string
	if sort == "" {
		sql = fmt.Sprintf("select * from product where name like '%%%s%%' or type like '%%%s%%' limit %d, %d", prodKeyWords, prodKeyWords, offset, pageSize)
	} else {
		sql = fmt.Sprintf("select * from product where name like '%%%s%%' or type like '%%%s%%' order by %s limit %d, %d", prodKeyWords, prodKeyWords, sort, offset, pageSize)
	}

	// 执行sql语句，获取查询结果集
	rows, err := dao.ExecQueryRows(sql)
	if err != nil {
		return nil, err
	}

	// 结果集内容转换为结构体
	prods := make([]model.Product, 0)
	for rows.Next() {
		var prod model.Product
		err = rows.Scan(&prod.ID, &prod.Name, &prod.Type, &prod.Price, &prod.Number, &prod.Images, &prod.Sales)
		if err != nil {
			return nil, err
		}
		prods = append(prods, prod)
	}

	// 返回
	return prods, nil
}

// GetChineseParameterValue 获取url中的中文参数值
func GetChineseParameterValue(URL string, regexStr string) (string, error) {
	// 转码
	URL, err := url.QueryUnescape(URL)
	if err != nil {
		return "", nil
	}

	// 正则表达式提取子字符串
	compileRegex, err := regexp.Compile(regexStr)
	if err != nil {
		return "", nil
	}
	matchStr := compileRegex.FindStringSubmatch(URL)

	return matchStr[1], nil
}

// AddComment 添加一条评论到数据库
func AddComment(comment model.Comment) (*model.Comment, error) {
	// 创建sql语句，忽略"ID"字段，实现在数据库中自增
	sql := CreateInsert(comment, "ID")

	// 执行sql语句将记录插入到数据库
	result, err := dao.ExecInsert(sql)

	// 取得数据库中自增字段"id"的值，把id(int 64)转为int
	idInt64, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	str := strconv.FormatInt(idInt64, 10)
	idInt, err := strconv.Atoi(str)
	if err != nil {
		return nil, err
	}

	// 补充comment内容
	comment.ID = idInt

	return &comment, nil
}

// CreateInsert 生成sql插入语句
func CreateInsert(v interface{}, passField string) string {
	// 获取结构体的值
	sValue := reflect.ValueOf(v)
	if sValue.Kind() == reflect.Ptr {
		sValue = sValue.Elem()
	}

	// 获取结构体的类型
	sType := sValue.Type()

	// 获取数据库表名(结构体名的小写)
	tableName := strings.ToLower(sType.Name())

	// 获取插入的字段数量
	filedNumber := sValue.NumField()

	// 生成sql语句
	isFirst := false
	sql := "insert into"
	sqlFiled := tableName + "("
	sqlValue := "value("
	for i := 0; i < filedNumber; i++ {
		curField := sType.Field(i).Name
		curValue := sValue.Field(i)

		// 生成语句时跳过指定字段
		if curField == passField {
			isFirst = true
			continue
		}

		switch curValue.Kind() {
		case reflect.String:
			if i == 0 || isFirst {
				sqlFiled += fmt.Sprintf("%s", curField)
				sqlValue += fmt.Sprintf("'%s'", curValue.String())
			} else {
				sqlFiled += fmt.Sprintf(", %s", curField)
				sqlValue += fmt.Sprintf(", '%s'", curValue.String())
			}
			isFirst = false
		case reflect.Int:
			if i == 0 || isFirst {
				sqlFiled += fmt.Sprintf("%s", curField)
				sqlValue += fmt.Sprintf("%d", curValue.Int())
			} else {
				sqlFiled += fmt.Sprintf(", %s", curField)
				sqlValue += fmt.Sprintf(", %d", curValue.Int())
			}
			isFirst = false
		case reflect.Float32, reflect.Float64:
			if i == 0 || isFirst {
				sqlFiled += fmt.Sprintf("%s", curField)
				sqlValue += fmt.Sprintf("%f", curValue.Float())
			} else {
				sqlFiled += fmt.Sprintf(", %s", curField)
				sqlValue += fmt.Sprintf(", %f", curValue.Float())
			}
			isFirst = false
		}
	}
	sqlFiled += ")"
	sqlValue += ")"
	sql = sql + " " + sqlFiled + " " + sqlValue + ";"

	return sql
}

// GetProdDetailInfo 获取某个商品的详细信息
func GetProdDetailInfo(prodID int) (*model.Product, error) {
	// 创建sql查询语句，分别为查询查询商品信息和查询评论信息
	sql1 := fmt.Sprintf("select * from product where id=%d;", prodID)
	sql2 := fmt.Sprintf("select * from comment where prodID=%d;", prodID)

	// 执行sql语句获取结果集
	row := dao.ExecQueryRow(sql1)
	rows, err := dao.ExecQueryRows(sql2)
	if err != nil {
		return nil, err
	}

	// 结果集转换为结构体
	var prod model.Product
	err = row.Scan(&prod.ID, &prod.Name, &prod.Type, &prod.Price, &prod.Number, &prod.Images, &prod.Sales)
	if err != nil {
		return nil, err
	}

	comments := make([]model.Comment, 0)
	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(&comment.ID, &comment.ProdID, &comment.ParentID, &comment.Content, &comment.UserName, &comment.Images)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	prod.Comments = comments

	// 返回
	return &prod, nil
}

// UpdateCart 更新购物车
func UpdateCart(userName string, prodID int) (*model.Cart, error) {
	// 查询该商品是否已经在购物车中，如果是则数量加1，否则插入新记录
	sqlStr := fmt.Sprintf("select * from cart where userName='%s' and prodID=%d;", userName, prodID)
	rows, err := dao.ExecQueryRows(sqlStr)
	if err != nil {
		return nil, err
	}

	length := 0
	for rows.Next() {
		length++
	}

	if length == 0 { //不存在
		cart := model.Cart{
			UserName: userName,
			ProdID:   prodID,
			Number:   1,
		}
		sqlStr = CreateInsert(cart, "")
	} else { //存在
		sqlStr = fmt.Sprintf("update cart set number=number+1 where userName='%s' and prodID=%d;", userName, prodID)
	}

	// 执行sql语句
	_, err = dao.ExecInsert(sqlStr)
	if err != nil {
		return nil, err
	}

	// 获取该记录的数据
	var record model.Cart
	sqlStr = fmt.Sprintf("select * from cart where userName='%s' and prodID=%d;", userName, prodID)
	row := dao.ExecQueryRow(sqlStr)
	err = row.Scan(&record.UserName, &record.ProdID, &record.Number)
	if err != nil {
		return nil, err
	}

	// 返回
	return &record, nil
}
