package controllers

import (
	"../models"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	this.Data["IsHome"] = true
	this.TplName = "home.html"

	this.Data["IsLogin"] = checkAccount(this.Ctx)

	topics,err := models.GetAllTopic(true)
	if err != nil {
		panic(err)
	}else {
		this.Data["Topics"] = topics
	}

	var errs error
	this.Data["Categories"],errs = models.GetAllCategories()
	if err != nil{
		beego.Error(errs)
	}
}
