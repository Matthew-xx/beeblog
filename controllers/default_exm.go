package controllers

import (
	"github.com/astaxie/beego"
)

type mainController struct {
	beego.Controller
}

func (c *MainController) get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"

	c.Data["TrueCond"] = true
	c.Data["FalseCond"] = false

	type U struct {
		Name string
		Age int64
		Sex string
	}
	user := &U{"mark",25,"man"}
	c.Data["Us"] = user   //with打印

	nums := []int{1,2,4,5,7}
	c.Data["Nums"] = nums  //range循环打印

	c.Data["Tmd"] = "burucesdhk"   //模板变量

}
