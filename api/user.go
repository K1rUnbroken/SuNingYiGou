package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"su-ning-yi-gou/model"
	"su-ning-yi-gou/service"
	"su-ning-yi-gou/util"
)

// Register 用户注册
func Register(c *gin.Context) {
	// 绑定表单参数到结构体
	var user model.User
	err := c.Bind(&user)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 用户名是否重复
	err, isRepeat := service.IsRepeat(user.UserName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if isRepeat {
		util.RespClientError(c, "user exists!")
	}

	// 添加该用户账号信息
	err = service.AddUser(user.UserName, user.Password)
	if err != nil {
		fmt.Println(err)
		return
	}

	util.RespSuccess(c, "register success!", user)
}

// Login 用户登录
func Login(c *gin.Context) {
	var user model.User
	err := c.Bind(&user)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 检查用户是否存在
	err, isExists := service.IsExists(user.UserName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !isExists {
		util.RespClientError(c, "user doesn't exists!")
		return
	}

	// 检查密码是否正确
	err, isRight := service.IsRight(user.UserName, user.Password)
	if !isRight {
		util.RespClientError(c, "wrong password!")
		return
	}

	// 生成jwt
	jwt, err := service.GenJWT(user.UserName)
	if err != nil {
		fmt.Println(err)
		return
	}

	user.JWT = jwt

	util.RespSuccess(c, "login success!", user)
}

// OutLogin 退出登录
func OutLogin(c *gin.Context) {

}

// SearchProduct 搜索商品
func SearchProduct(c *gin.Context) {
	// 获取url中的普通的参数值
	pageSize1 := c.Query("pageSize")                 //每页显示的商品数量
	pageNumber1 := c.DefaultQuery("pageNumber", "1") //显示第几页，默认为第1页
	sort := c.DefaultQuery("sort", "")               //排序依据

	// 获取url中的值为中文的参数值（刚开始不知道什么原因不能直接获取中文值参数，所以做了单独处理。后来发现又可以直接获取，但还是保留了这个方法）
	URL := fmt.Sprint(c.Request.URL)
	prodKeyWords, err := service.GetChineseParameterValue(URL, ".*prodName=(.*?)&.*")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 参数类型转换
	pageSize2, err := strconv.Atoi(pageSize1)
	if err != nil {
		fmt.Println(err)
		return
	}
	pageNumber2, err := strconv.Atoi(pageNumber1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 查询结果集
	offset := (pageNumber2 - 1) * pageSize2 // 偏移量，即从查询结果集中的第几条记录开始显示
	prods, err := service.GetProdByNameAndType(prodKeyWords, offset, pageSize2, sort)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 返回json
	util.RespSuccess(c, "success!", prods)
}

// ShowProductByType 分类显示商品
func ShowProductByType(c *gin.Context) {
	// 获取表单参数

	// 查询结果集

}

// ProductDetail 商品详细信息
func ProductDetail(c *gin.Context) {
	// 获取url中的参数
	id1 := c.Query("prodID")

	// 参数类型转换
	id2, err := strconv.Atoi(id1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 查询结果集
	prod, err := service.GetProdDetailInfo(id2)

	// 返回json
	util.RespSuccess(c, "success!", *prod)
}

// PostComment 发布一条评论
func PostComment(c *gin.Context) {
	//表单部分参数绑定到结构体
	var comment model.Comment
	err := c.Bind(&comment)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取表单上传的多张图片(images字段未绑定，单独处理)
	form, err := c.MultipartForm()
	images := form.File["images"]
	if err != nil {
		fmt.Println(err)
		return
	}

	// 保存图片到服务器指定路径,得到所有图片的url切片
	savePath := "images/product/"
	urlSli, err := service.UploadImagesToServer(savePath, images, c)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 以空格" "将url切片中所有url拼接
	var urls string
	number := len(urlSli)
	for i := 0; i < number; i++ {
		if i+1 < number {
			urls = urls + urlSli[i] + " "
		} else {
			urls += urlSli[i]
		}
	}

	// 将comment中的Images属性值改为拼接后的字符串
	comment.Images = urls

	// 存储到数据库，返回该评论的完整信息
	newComment, err := service.AddComment(comment)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 返回json
	util.RespSuccess(c, "post comment success!", *newComment)
}

// AddCart 商品加入到购物车
func AddCart(c *gin.Context) {
	// 获取url中的参数
	userName := c.Query("userName")
	prodID1 := c.Query("prodID")

	// 参数类型转换
	prodID2, err := strconv.Atoi(prodID1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 更新购物车
	info, err := service.UpdateCart(userName, prodID2)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 返回
	util.RespSuccess(c, "add cart success!", *info)
}
