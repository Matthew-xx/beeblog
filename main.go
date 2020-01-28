package main

import (
	_ "beeblog/routers"
	"github.com/astaxie/beego"
	"./models"
	"./controllers"
	"github.com/astaxie/beego/orm"
	"os"
)

func init()  {
	models.RegisterDB() //初始化数据库
}

func main() {
	orm.Debug = true  //方便调试看数据库是否创建
	orm.RunSyncdb("default",false,true) //自动建表

	//（管理）登录操作
	beego.Router("/",&controllers.MainController{})
	beego.Router("/home" + "",&controllers.MainController{})
	beego.Router("/login",&controllers.LoginController{})
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/category",&controllers.CategoryController{})
	beego.Router("/topic",&controllers.TopicController{})
	beego.Router("/reply",&controllers.ReplyController{})
	beego.Router("/reply/add",&controllers.ReplyController{},"post:Add")
	beego.Router("/reply/delete",&controllers.ReplyController{},"get:Delete")

	//附件处理
	//创建附件目录
	os.Mkdir("attachment",os.ModePerm)

	//作为静态文件处理
	//beego.SetStaticPath("/attachment","attachment")

	//作为单独一个控制器来处理
	beego.Router("/attachment/:all",&controllers.AttachController{})


	beego.Run()
}

