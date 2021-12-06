package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//定义一个结构体
type User struct {
	Id int
	Name string
	PassWord string
	//Pass_Word
	Articles []*Article `orm:"reverse(many)"`
}

type Article struct {
	Id int `orm:"pk;auto"`
	ArtiName string `orm:"size(20)"`
	Atime time.Time `orm:"auto_now"`
	Acount int `orm:"default(0);null"`
	Acontent string `orm:"size(500)"`
	Aimg string  `orm:"size(100);null"`

	ArticleType *ArticleType `orm:"rel(fk)"` //fk外键
	Users []*User `orm:"rel(m2m)"`
}

//类型表
type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`

	Articles []*Article `orm:"reverse(many)"` //rel和reverse成对存在
}

func init()  {
	orm.RegisterDataBase("default", "mysql", "root:root2@tcp(127.0.0.1:3306)/test?charset=utf8")
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	orm.RunSyncdb("default", false, true)
}