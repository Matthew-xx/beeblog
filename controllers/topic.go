package controllers

import (
	"../models"
	"github.com/astaxie/beego"
	"path"
	"strings"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController) Get()  {
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	this.Data["IsTopic"] =true
	this.TplName = "topic.html"
	topics,err := models.GetAllTopic("","",false)
	if err != nil {
		beego.Error(err)
	}else {
		this.Data["Topics"] = topics
	}
}

func (this *TopicController) Post()  {

	if !checkAccount(this.Ctx)  {
		this.Redirect("/login",302)
		return
	}

	//解析表单
	tid := this.Input().Get("tid")
	title := this.Input().Get("title")
	content := this.Input().Get("content")
	label := this.Input().Get("label")
	category := this.Input().Get("category")

	//判断用户是否上传附件
	_,fh,err := this.GetFile("attachment")
	if err != nil{
		beego.Error(err)
	}

	var attachment string
	if fh != nil {
		//保存附件
		attachment = fh.Filename
		beego.Info(attachment)
		err = this.SaveToFile("attachment",path.Join("attachment",attachment))
		// 第二参数（可相对路径)里面参数（文件夹，文件）
		if err != nil {
			beego.Error(err)
		}
	}

	if len(tid) == 0 {
		err = models.AddTopic(title,category,label,content,attachment)
	}else {
		err = models.ModifyTopic(tid,title,category,label,content,attachment)
	}

	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/topic",302)
}

func (this *TopicController) Add()  {
	this.TplName = "topic_add.html"
}

func (this *TopicController) View() {
	this.TplName = "topic_view.html"

	tid := this.Input().Get("tid")
	topic,err := models.GetTopic(tid)
	//topic,err := models.GetTopic(this.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
		this.Redirect("/home",302)
		return
	}

	this.Data["Topic"]=topic
	this.Data["Labels"]=strings.Split(topic.Labels," ")
	this.Data["Tid"] = tid

	replies,err := models.GetAllReplies(tid)
	if err != nil {
		beego.Error(err)
		return //不影响其他查看步骤
	}

	this.Data["Replies"] = replies
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}

func (this *TopicController) Modify() {
	this.TplName = "topic_modify.html"

	tid := this.Input().Get("tid")
	topic,err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/home",302)
		return
	}

	this.Data["Topic"]=topic
	this.Data["Tid"] = tid
}

func (this *TopicController) Delete() {
	if !checkAccount(this.Ctx)  {
		this.Redirect("/login",302)
		return
	}

	err := models.DeleteTopic(this.Input().Get("tid"))
	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/home",302)
}

