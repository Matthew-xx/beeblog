package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	_ "github.com/go-sql-driver/mysql"
)


const(
	_DB_NAME   = "beeblog"
	_DRIVER_NAME ="mysql"
)
type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
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
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}

func RegisterDB()  {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	// set default database
	orm.RegisterDataBase("default", "mysql", "root:666666@tcp(127.0.0.1:3306)/beeblog?charset=utf8",30)

	orm.RegisterModel(new(Category),new(Topic))
}

//添加文章
func AddCategory(name string) error {
	o := orm.NewOrm()

	cate := &Category{Title: name}
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


