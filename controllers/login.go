package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	isExit := this.Input().Get("exit")
	if isExit == "true" { //重新设置cookie
		this.Ctx.SetCookie("uname","",-1,"/") // -1 立即删除的效果
		this.Ctx.SetCookie("pwd","",-1,"/")
		this.Redirect("/",301)
		return
	}
	this.TplName = "login.html"
}

//一般是post登录
func (this *LoginController) Post()  {
	uname := this.Input().Get("uname")
	pwd := this.Input().Get("pwd")
	autologin := this.Input().Get("autologin") == "on"

	if beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("pwd")==pwd {
		maxAge := 0
		if autologin {
			maxAge = 1<<31 - 1
		} //账号保存时间（关闭浏览器
		this.Ctx.SetCookie("uname",uname,maxAge,"/")
		this.Ctx.SetCookie("pwd",pwd,maxAge,"/")
	}
	this.Redirect("/",301) //重定向
	return
}

//读取cookie判断验证用户是否登录
func checkAccount(ctx *context.Context) bool {
    ck,err := ctx.Request.Cookie("uname")
	if err != nil {
		return false
	}
    uname := ck.Value

    ck,err = ctx.Request.Cookie("pwd")
	if err != nil {
		return false
	}
    pwd := ck.Value

    return beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("pwd") == pwd
}
