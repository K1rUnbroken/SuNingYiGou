package api

// VerifyJWT 校验jwt（jwt放在Header的Authorization中，并使用Bearer开头）
//func VerifyJWT(c *gin.Context) {
//	auth := c.Request.Header.Get("")
//
//	// 检查auth是否为空
//	if auth == "" {
//		util.ResponseSuccess(c, "empty authorization!")
//		c.Abort()
//		return
//	}
//
//	// 检查auth格式是否正确
//	parts := strings.SplitN(auth, " ", 2)
//	if !(len(parts) == 2 && parts[0] == "Bearer") {
//		util.ResponseSuccess(c, "wrong authorization format!")
//		c.Abort()
//		return
//	}
//
//	// 解析获取到的jwt
//	claims, err := service.ParseJWT(parts[1])
//	if err != nil {
//		util.ResponseSuccess(c, "invalid jwt!")
//		c.Abort()
//		return
//	}
//
//	// 将当前请求的userName保存到上下文中
//	c.Set("userName", claims.UserName)
//
//	c.Next()
//}
