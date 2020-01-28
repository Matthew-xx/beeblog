package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path"
	"strconv"
	"strings"
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
	Category        string
	Labels          string
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

type Comment struct {
	Id int64
	Tid int64
	Name string
	Content string `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

func RegisterDB()  {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	// set default database
	orm.RegisterDataBase("default", "mysql", "root:666666@tcp(127.0.0.1:3306)/beeblog?charset=utf8",30)

	orm.RegisterModel(new(Category),new(Topic),new(Comment))  //每添加新表需先注册
}

func AddTopic(title,category,label,content,attachment string) error {
	//新增标签，处理标签（实现简单搜索功能，
	label = "$"+strings.Join(strings.Split(label," "),"#$") + "#"  //以空格作为多个标签的分隔符
	//“ $label# " 存入数据库形式。关键字搜索，模糊搜索
	//如传入“bee orm“ 那么strings处理后就是  bee#$orm 最后是 $bee#$orm#

	o := orm.NewOrm()
	topic := &Topic{
		Title:title,
		Content:content,
		Labels:label,
		Attachment:attachment,
		Category:category,
		Created:time.Now(),
		Updated:time.Now(),
	}

	_,err := o.Insert(topic)
	if err != nil{
		beego.Error(err)
	}
	//更新分类统计
	cate := new(Category)  //获取分类对象
	qs := o.QueryTable("category")
	err = qs.Filter("title",category).One(cate)
	if err == nil {
		//如果不存在,简单忽略更新操作
		cate.TopicCount++
		_,err = o.Update(cate)
	}

	return nil
}

func AddReply(tid,nickname,content string) error {
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return err
	}

	reply := &Comment{
		Tid:tidNum,
		Name:nickname,
		Content:content,
		Created:time.Now(),
	}

	o := orm.NewOrm()
	_,err = o.Insert(reply)
	if err != nil {
		return err
	}

	topic := &Topic{Id:tidNum}
	if o.Read(topic) == nil{
		topic.ReplyTime = time.Now()
		topic.ReplyCount++
		_,err = o.Update(topic)
	}
	return err
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

func GetAllReplies(tid string) (replies []*Comment,err error) {
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return nil,err
	}
	replies = make([]*Comment,0)

	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_,err = qs.Filter("tid",tidNum).All(&replies)

	return replies,err
}

func GetAllTopic(cate ,label string ,isDesc bool) ([]*Topic,error) {
	//倒序
	o := orm.NewOrm()
	topics := make([]*Topic,0)

	qs := o.QueryTable("topic")


	var err error
	if isDesc{
		if len(cate) > 0 {
			qs = qs.Filter("category",cate)  //过滤，保存过滤后的对象
		} //看是否有分类
		if len(label) > 0 {
			qs = qs.Filter("labels__contains","$"+label+"#")  //查找包含"$"+label+"#"的labels
		}//看是否有标签
		_,err = qs.OrderBy("-created").All(&topics)  //链式操作，所有结果读取完后一次性返回
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

	//将数据库存储的形式转换成输入形式
	topic.Labels = strings.Replace(strings.Replace(topic.Labels,"#"," ",-1),"$","",-1)
	return topic,err
}

func ModifyTopic(tid,title,category,label,content,attachment string) error {
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return err
	}

	label = "$"+strings.Join(strings.Split(label," "),"#$") + "#"

	var oldCate,oldAttach string
	o := orm.NewOrm()
	topic := &Topic{Id:tidNum}

	if o.Read(topic) == nil {
		oldCate = topic.Category  //先取得旧的分类名称
		oldAttach = topic.Attachment
		topic.Title = title
		topic.Category = category  //再将新的分类名称赋值
		topic.Attachment = attachment
		topic.Labels = label
		topic.Content=content
		topic.Updated=time.Now()
		o.Update(topic)
		if err != nil {
			return err
		}
	}

	//更新分类统计
	if len(oldCate) > 0{
		cate := new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("title",oldCate).One(cate)  //找到旧的分类
		if err == nil {
			cate.TopicCount--  //若能找到旧的分类则减去计数
			_,err = o.Update(cate)
		}
	}  //更新旧的分类统计

	//删除旧的附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment",oldAttach))
	}

	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title",category).One(cate)
	if err == nil {
		cate.TopicCount++  //若能找到旧的分类则减去计数
		_,err = o.Update(cate)
	}//更新新的分类统计

	return err
}

func DeleteTopic(tid string) error {
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return err
	}

	var oldCate string
	o := orm.NewOrm()
	topic := &Topic{Id:tidNum}
	//先获取文章分类
	if o.Read(topic) == nil {
		oldCate = topic.Category
		_,err = o.Delete(topic) //已经存到临时变量里，删除无影响
		if err != nil {
			return err
		}
	}

	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title",oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_,err = o.Update(cate)
		}
	}  //文章删除后统计减去

	return err
}

func DeleteReply(rid string) error {
	ridNum,err := strconv.ParseInt(rid,10,64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()

	//先存储到临时变量中
	var tidNum int64

	reply := &Comment{Id:ridNum}
	if o.Read(reply) == nil{
		//返回错误为空即能读取出来，再进行下一步
		tidNum = reply.Id
		_,err = o.Delete(reply)  //查找id相等的记录再删除（先初始化一个对象
		if err != nil {
			return err
		}
	}

	replies := make([]*Comment,0)  //使用字典实现精确统计，所有相关回复都取出再取其长度（不使用reply--）以及方便正确获取最后回复时间
	qs := o.QueryTable("comment")
	_,err = qs.Filter("tid",tidNum).OrderBy("-created").All(&replies)
	if err != nil {
		return err
	}

	topic := &Topic{Id:tidNum}
	if o.Read(topic) == nil{
		topic.ReplyTime = replies[0].Created   //获取最后回复的创建时间
		topic.ReplyCount = int64(len(replies))
		_,err = o.Update(topic)
	}
	return err
}
