package service

import (
	"fmt"
	"ginchat/utils"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	w := c.Writer
	req := c.Request
	srcFile, head, err := req.FormFile("file")
	if err != nil {
		utils.ResFail(w, err.Error())
	}
	suffix := ".png"
	ofilName := head.Filename
	tem := strings.Split(ofilName, ".png")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFile, err := os.Create("asset/upload/" + fileName)
	if err != nil {
		utils.ResFail(w, err.Error())
	}
	_, errCopy := io.Copy(dstFile, srcFile)
	if errCopy != nil {
		utils.ResFail(w, err.Error())
	}
	url := "./asset/upload/" + fileName
	utils.ResOK(w, url, "发送图片成功")
}
