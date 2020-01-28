package controllers

import (
	"github.com/astaxie/beego"
	"io"
	"net/url"
	"os"
)

type AttachController struct {
	beego.Controller
}

func (this *AttachController) Get() {
	filePath,err := url.QueryUnescape(this.Ctx.Request.RequestURI[1:])
	// 截取后的路径是  /attachment/go-blog功能.png  需取得第一个“/"
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}

	f,err := os.Open(filePath)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}
	defer f.Close()

	_,err = io.Copy(this.Ctx.ResponseWriter,f)
}
