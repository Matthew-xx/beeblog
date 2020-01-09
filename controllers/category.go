package controllers

import "github.com/astaxie/beego"

type CategoryController struct {
	beego.Controller
}

func (this *CategoryController) Get()  {
	this.Data["IsCategory"] = true
	this.TplName = "category.html"

}
