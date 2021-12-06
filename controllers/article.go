package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"myproject/models"
	"path"
	"strconv"
)

type ArticleController struct {
	beego.Controller
}

//展示文章列表页
func(this*ArticleController)ShowArticleList(){
	//session判断
	userName := this.GetSession("userName")
	if userName == nil{
		this.Redirect("/",302)
		return
	}

	//获取数据
	//高级查询
	//指定表
	o := orm.NewOrm()
	qs := o.QueryTable("Article")//queryseter
	//var articles []models.Article

	//查询总记录数
	typeName := this.GetString("select")
	var count int64

	//获取总页数
	pageSize := 2

	//获取页码
	pageIndex,err:= this.GetInt("pageIndex")
	if err != nil{
		pageIndex = 1
	}

	//获取数据
	//作用就是获取数据库部分数据,第一个参数，获取几条,第二个参数，从那条数据开始获取,返回值还是querySeter
	//起始位置计算
	start := (pageIndex - 1)*pageSize

	//qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)

	if typeName == ""{
		count,_ = qs.Count()
	}else{
		count,_ = qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
	}
	pageCount := math.Ceil(float64(count) / float64(pageSize))

	//获取文章类型
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types

	//根据选中的类型查询相应类型文章
	var articleswithtype []models.Article
	beego.Info(typeName)
	if typeName == ""{
		//qs.Limit(pageSize,start).All(&articles)
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articleswithtype)

	}else {
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articleswithtype)
	}

	//传递数据
	this.Data["userName"] = userName.(string)
	this.Data["typeName"] = typeName
	this.Data["pageIndex"] = pageIndex
	this.Data["pageCount"] = int(pageCount)
	this.Data["count"] = count
	this.Data["articles"] = articleswithtype

	//指定试图布局
	this.Layout = "layout.html"
	this.TplName = "index.html"
}

//处理下拉框改变发的请求
func (this *ArticleController) HandleSelect(){
	//1.接收数据
	typeName := this.GetString("select")
	//beego.Info(typeName)
	//2.处理数据
	if typeName == ""{
		beego.Info("下拉框传递数据失败")
		return
	}

	//3.查询数据
	o := orm.NewOrm()
	var articles []models.Article
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)

	this.Ctx.WriteString(typeName)
}

func (this *ArticleController) ShowContent(){
	//1.获取id
	id := this.GetString("id")
	//beego.Info(id)

	//2.查询
	o := orm.NewOrm()
	id2 , _ := strconv.Atoi(id)
	article := models.Article{Id:id2}
	err := o.Read(&article)
	if err != nil {
		beego.Info("查询数据为空")
		return
	}

	article.Acount += 1
	o.Update(&article)

	//3.展示到视图
	this.Data["article"] = article
	this.Layout = "layout.html"
	this.TplName = "content.html"
}

func (this * ArticleController) ShowAddArticle(){
	//查询类型数据，传递到视图中
	var types []models.ArticleType
	o := orm.NewOrm()
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types

	this.TplName = "add.html"
}

/**
	1. 拿数据
	2. 判断数据
	3. 插入数据
	4. 返回视图
 */
//获取添加文章数据
func(this*ArticleController)HandleAddArticle(){
	//1.获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")

	//2校验数据
	if articleName == "" || content == ""{
		this.Data["errmsg"] = "添加数据不完整"
		this.TplName = "add.html"
		return
	}

	//处理文件上传
	file , head,err:=this.GetFile("filename")
	defer file.Close()
	if err != nil{
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return
	}


	//1.文件大小
	if head.Size > 5000000{
		this.Data["errmsg"] = "文件太大，请重新上传"
		this.TplName = "add.html"
		return
	}

	//2.文件格式
	//a.jpg
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		this.Data["errmsg"] = "文件格式错误。请重新上传"
		this.TplName = "add.html"
		return
	}

	//3.防止重名
	//fileName := time.Now().Format("2006-01-02-15:04:05")
	//存储
	this.SaveToFile("filename","static/upload/"+ head.Filename)

	//3.处理数据
	//插入操作
	o := orm.NewOrm()

	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	var article models.Article
	article.ArtiName = articleName
	article.Acontent = content
	article.Aimg = "static/upload/"+ head.Filename
	//给文章添加类型
	//获取类型数据
	typeName := this.GetString("select")
	//根据名称查询类型
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType,"TypeName")

	article.ArticleType = &articleType

	o.Insert(&article)

	this.Data["types"] = types
	//4.返回页面
	this.Redirect("/Article/showArticleList",302)
}

//1. URL传值
//2. Delete
func (this *ArticleController) HandleDelete(){

	id, _ := this.GetInt("id")
	//1.orm对象
	o := orm.NewOrm()
	//2.删除对象
	article := models.Article{Id: id}
	//3.删除
	o.Delete(&article)

	this.Redirect("/Article/showArticleList", 302)
}

func (this *ArticleController) ShowUpdate(){
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	article := models.Article{Id: id}
	err := o.Read(&article)
	if err != nil {
		beego.Info("查询数据为空")
		return
	}
	this.Data["article"] = article

	this.TplName = "update.html"
}

func (this *ArticleController) HandleUpdate(){
	id, err := this.GetInt("id")
	articleName := this.GetString("articleName")
	artiContent := this.GetString("content")

	if err != nil {
		beego.Info("Update无法获取id")
		return
	}

	beego.Info(id)

	file , head, err := this.GetFile("filename")
	beego.Info(err)

	defer file.Close()
	if err != nil{
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return
	}
	beego.Info(id)
	//1.文件大小
	if head.Size > 5000000{
		this.Data["errmsg"] = "文件太大，请重新上传"
		this.TplName = "add.html"
		return
	}
	//2.文件格式
	//a.jpg
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		this.Data["errmsg"] = "文件格式错误。请重新上传"
		//this.TplName = "add.html"
		return
	}
	//3.防止重名
	//fileName := time.Now().Format("2006-01-02-15:04:05")
	//存储
	this.SaveToFile("filename","static/upload/"+ head.Filename)
	beego.Info(id)
	//3.处理数据
	//更新操作
	o := orm.NewOrm()
	article := models.Article{Id: id}

	err = o.Read(&article)
	if err != nil{
		beego.Info("更新的文章不存在")
		return
	}

	article.ArtiName = articleName
	article.Acontent = artiContent
	article.Aimg = "static/upload/"+ head.Filename
	o.Update(&article)

	//4.返回页面
	this.Redirect("/Article/showArticleList",302)
}

func (this *ArticleController) ShowAddType(){
	//1. 读取类型表，显示数据
	o := orm.NewOrm()
	var types []models.ArticleType
	//查询
	o.QueryTable("article_type").All(&types)

	//beego.Info(types)
	this.Data["types"] = types

	this.TplName = "addType.html"
}
//添加类型业务
func (this *ArticleController) HandleAddType(){
	//1 获取数据
	typename := this.GetString("typeName")
	//2 判断数据
	if typename == ""{
		beego.Info("添加类型数据不能为空")
		return
	}
	//3 执行插入操作
	o := orm.NewOrm()
	var artiType models.ArticleType
	artiType.TypeName =typename
	_, err :=o.Insert(&artiType)
	if err != nil {
		beego.Info("插入失败")
	}
	//4 展示视图
	this.Redirect("/Article/addArticleType", 302)
}