package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// InitRouter 初始化路由
func InitRouter() {
	r := gin.Default()

	// 设置为静态文件夹
	r.StaticFS("/images/product", http.Dir("./images/product"))
	r.StaticFS("/images/comment", http.Dir("./images/comment"))

	// 用户
	userGroup := r.Group("/user")
	{
		userGroup.PUT("/register", Register) //注册
		userGroup.PUT("/login", Login)       //登录
		userGroup.PUT("/outLogin", OutLogin) //退出登录
		userGroup.POST("/addCart", AddCart) //商品加入到购物车

		prodGroup := userGroup.Group("/product")
		{
			prodGroup.GET("/search", SearchProduct)     //搜索商品(商城首页搜索栏)
			prodGroup.GET("/show", ShowProductByType)   //分类显示商品(商城首页分类栏)
			prodGroup.GET("/detail", ProductDetail)     //商品详细
			prodGroup.POST("/postComment", PostComment) //发布评论
		}

	}

	// 卖家
	sellerGroup := r.Group("/seller")
	{
		sellerGroup.PUT("/addProduct", AddProd) //卖家上架商品
	}

	r.Run()
}
