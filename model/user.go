package model

import (
	"github.com/dgrijalva/jwt-go"
)

// User 用户
type User struct {
	UserName string `form:"userName" json:"userName"`
	Password string `form:"password" json:"password"`
	JWT      string `form:"-" json:"jwt"`
}

// Claims jwt的payload
type Claims struct {
	UserName string `json:"userName"`
	jwt.StandardClaims
}

// Product 商品
type Product struct {
	ID       int       `form:"-" json:"id"`
	Name     string    `form:"name" json:"name"`
	Type     string    `form:"type" json:"type"`
	Price    float64   `form:"price" json:"price"`
	Number   int       `form:"number" json:"number"`
	Images   string    `form:"-" json:"imagesURL"`
	Sales    int       `form:"-" json:"sales"`
	Comments []Comment `form:"-" json:"comments"`
}

// Comment 商品评论
type Comment struct {
	ID       int    `form:"-" json:"id"`
	ProdID   int    `form:"prodID" json:"prodID"`
	ParentID int    `form:"parentID" json:"parentID"`
	Content  string `form:"content" json:"content"`
	UserName string `form:"userName" json:"userName"`
	Images   string `form:"-" json:"imagesURL"`
}

//// ResultSet 查询商品的结果集
//type ResultSet []map[string]any

// Cart 购物车
type Cart struct {
	UserName string `json:"userName"`
	ProdID   int    `json:"prodID"`
	Number   int    `json:"number"`
}
