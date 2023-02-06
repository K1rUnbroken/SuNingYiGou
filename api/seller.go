package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"su-ning-yi-gou/model"
	"su-ning-yi-gou/service"
	"su-ning-yi-gou/util"
)

// AddProd 卖家上架新商品
func AddProd(c *gin.Context) {
	// 绑定表单部分参数
	var prod model.Product
	err := c.Bind(&prod)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取表单上传的图片
	form, err := c.MultipartForm()
	images := form.File["images"]
	if err != nil {
		fmt.Println(err)
		return
	}

	// 检查商品名是否重复

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

	// 补充"Images"和"Sales"字段的值
	prod.Images = urls
	prod.Sales = 0

	// 添加商品到数据库
	newProd, err := service.AddProduct(prod)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 返回json
	util.RespSuccess(c, "add product success!", *newProd)
}
