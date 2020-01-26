package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)


const(
	_DB_NAME   = "beeblog"
	_DRIVER_NAME ="mysql"
)
type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index;null"`
	Views           int64     `orm:"index;null"`
	TopicTime       time.Time `orm:"index;null"`
	TopicCount      int64
	TopicLastUserId int64
}

//反射的时候可以得到tag作为一个说明得到`orm:"index"`,表示只有orm可读，设置值为index

type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Content         string `orm:size(5000)`
	Attachment      string
	Created         time.Time `orm:"index;null"`
	Updated         time.Time `orm:"index;null"`
	Views           int64     `orm:"index;null"`
	Author          string
	ReplyTime       time.Time `orm:"index;null"`
	ReplyCount      int64
	ReplyLastUserId int64
}

func RegisterDB()  {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	// set default database
	orm.RegisterDataBase("default", "mysql", "root:666666@tcp(127.0.0.1:3306)/beeblog?charset=utf8",30)

	orm.RegisterModel(new(Category),new(Topic))
}

func AddTopic(title,content string) error {
	o := orm.NewOrm()
	topic := &Topic{
		Title:title,
		Content:content,
		Created:time.Now(),
		Updated:time.Now(),
	}

	_,err := o.Insert(topic)
	if err != nil{
		beego.Error(err)
	}

	return nil
}

//添加文章
func AddCategory(name string) error {
	o := orm.NewOrm()
	tt := time.Now()

	cate := &Category{Title: name,Created: tt}
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate) //使用one获取单个对象

	if err == nil {
		return err
	}

	_, err = o.Insert(cate)
	if err != nil {  //插入失败
		return err
	}
	return nil
}

func GetAllTopic(isDesc bool) ([]*Topic,error) {
	//倒序
	o := orm.NewOrm()
	topics := make([]*Topic,0)

	qs := o.QueryTable("topic")


	var err error
	if isDesc{
		_,err = qs.OrderBy("-created").All(&topics)
	}else {
		_,err = qs.All(&topics)
	}

	return topics,err
}

//获取
func GetAllCategories() ([]*Category,error)  {
	o := orm.NewOrm()
	cates := make([]*Category,0)
	qs := o.QueryTable("category")

	_,err := qs.All(&cates)

	return cates,err
}

func DelCategories(id string) error {
	cid ,err := strconv.ParseInt(id,10,64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id:cid}
	_,err = o.Delete(cate)

	return err
}

func GetTopic(tid string) (*Topic, error) {
	/*
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return nil,err
	}
*/
	o := orm.NewOrm()
	topic := new(Topic)

	qs := o.QueryTable("topic")
	err := qs.Filter("id",tid).One(topic)
	if err != nil {
		return nil,err
	}

	topic.Views++
	_,err = o.Update(topic)
	return topic,err
}

func ModifyTopic(tid,title,content string) error {
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	topic := &Topic{Id:tidNum}

	if o.Read(topic) == nil {
		topic.Title = title
		topic.Content=content
		topic.Updated=time.Now()
		o.Update(topic)
	}

	return err
}

func DeleteTopic(tid string) error {
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	topic := &Topic{Id:tidNum}
	_,err = o.Delete(topic)

	return err
}
