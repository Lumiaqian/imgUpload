package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		//fmt.Println(method)
		c.Header("Access-Control-Allow-Origin", "*")
		//c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X_Tk, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(cors())
	router.POST("/uploadHandler", UploadHandler)
	router.GET("/fileDownload", FileDownload)
	err := router.Run(":10086")
	if err != nil {
		return
	}
}

// UploadHandler 单张图片上传
func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("uploadfile")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"msg":     fmt.Sprintf("图片上传失败！%v", err),
		})
		log.Panicln("图片上传失败！", err)
		return
	}
	name := ""
	// 文件夹
	folder := "/app/data/images/"
	//如果没有filepath文件目录就创建一个
	if _, err := os.Stat(folder); err != nil {
		if !os.IsExist(err) {
			err := os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				return
			}
		}
	}
	// 域名
	prefix := "https://img.caoyuqian.cn/"
	if file != nil {
		// 获得文件名
		name = strings.TrimSpace(file.Filename)
		err := c.SaveUploadedFile(file, folder+name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"msg":     fmt.Sprintf("图片上传失败！%v", err),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "desc": "图片上传成功", "data": prefix + name})

}

func FileDownload(c *gin.Context) {
	filename, err := c.GetQuery("filename")
	if !err {
		c.String(400, "Success")
	}

	path := "./"
	path += filename
	fmt.Println(path)
	c.File(path)
}
