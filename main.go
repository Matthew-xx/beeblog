package main

import (
	_ "beeblog/routers"
	"github.com/astaxie/beego"
	"./models"
	"./controllers"
	"github.com/astaxie/beego/orm"
)

func init()  {
	models.RegisterDB() //初始化数据库
}

func main() {
	orm.Debug = true  //方便调试看数据库是否创建
	orm.RunSyncdb("default",false,true)

	//（管理）登录操作
	beego.Router("/",&controllers.MainController{})
	beego.Router("/home" +
		"",&controllers.MainController{})
	beego.Router("/login",&controllers.LoginController{})
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/category",&controllers.CategoryController{})
	beego.Router("/topic",&controllers.TopicController{})
	//beego.Router("/topic/add",&controllers.TopicController{})
	beego.Run()
}

