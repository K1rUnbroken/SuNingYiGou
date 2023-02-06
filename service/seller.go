package service

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path"
	"strconv"
	"strings"
	"su-ning-yi-gou/dao"
	"su-ning-yi-gou/model"
	"su-ning-yi-gou/util"
)

// AddProduct 新增商品
func AddProduct(prod model.Product) (*model.Product, error) {
	// 创建sql语句，忽略"ID"字段，实现在数据库中自增
	sql := CreateInsert(prod, "ID")

	// 执行sql语句将记录插入到数据库
	result, err := dao.ExecInsert(sql)
	if err != nil {
		return nil, err
	}

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

	// 还原comment内容
	prod.ID = idInt

	return &prod, nil
}

// UploadImagesToServer 下载图片保存在服务器，返回图片的url切片
func UploadImagesToServer(savePath string, images []*multipart.FileHeader, c *gin.Context) ([]string, error) {
	var urlSli []string // 所有图片的url切片

	// 每个图片保存到服务器，得到其url添加到切片中
	for _, image := range images {
		// 根据md5编码生成新的文件名
		imageName := func(name string) string {
			ext := path.Ext(name)
			fileName := strings.TrimSuffix(name, ext)
			fileName = util.EncodeMD5(fileName)

			return fileName + ext
		}(image.Filename)

		// 图片保存到服务器，其位置为src，得到url
		src := savePath + imageName
		err := c.SaveUploadedFile(image, src)
		if err != nil {
			return nil, err
		}
		url := "http://127.0.0.1:8080/" + src

		// url添加到切片中
		urlSli = append(urlSli, url)
	}

	return urlSli, nil
}
