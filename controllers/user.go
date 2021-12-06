package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"myproject/models"
)

type UserController struct {
	beego.Controller
}

//显示注册页面
func(this *UserController)ShowRegister(){
	this.TplName = "register.html"
}


//处理注册数据
func(this *UserController)HandlePost(){
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")

	//beego.Info(userName,pwd)
	//2.校验数据
	if userName == "" || pwd == ""{
		this.Data["errmsg"] = "注册数据不完整，请重新注册"
		beego.Info("注册数据不完整，请重新注册")
		this.TplName = "register.html"
		return
	}
	//3.操作数据
	//获取ORM对象
	o := orm.NewOrm()
	//获取插入对象
	user := models.User{}
	//给插入对象赋值
	user.Name = userName
	user.PassWord = pwd
	//插入
	_,err := o.Insert(&user)
	if err != nil{
		beego.Info("插入失败")
	}
	//返回结果

	//4.返回页面
	//this.Ctx.WriteString("注册成功")
	this.Redirect("/",302)
	//this.TplName = "login.html"
}


//展示登录页面
func(this *UserController) ShowLogin(){
	userName := this.Ctx.GetCookie("userName")
	data,_ := base64.StdEncoding.DecodeString(userName)
	if userName == ""{
		this.Data["userName"] = ""
		this.Data["check"] = ""
	}else{
		this.Data["userName"] = string(data)
		this.Data["check"] = "checked"
	}
	this.TplName = "login.html"
}


func(this*UserController) HandleLogin(){
	fmt.Println("HandleLogin")
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	beego.Info(userName)
	//2.校验数据
	if userName == "" || pwd == ""{
		this.Data["errmsg"] = "登录数据不完整"
		this.TplName = "login.html"
		return
	}

	fmt.Println("ORM")
	//3.操作数据
	//1。获取ORM对象
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	err := o.Read(&user,"Name")
	beego.Info(userName)
	if err != nil{
		this.Data["errmsg"] = "用户不存在"
		this.TplName = "login.html"
		return
	}
	if user.PassWord != pwd{
		this.Data["errmsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}


	//4.返回页面
	//this.Ctx.WriteString("")
	data := this.GetString("remember")
	beego.Info(data)

	if data == "on"{
		this.Ctx.SetCookie("userName",userName,100)
	}else {
		this.Ctx.SetCookie("userName",userName,-1)
	}


	this.SetSession("userName",userName)
	this.Redirect("/Article/showArticleList",302)

}

//退出登录
func(this *UserController) Logout(){
	//删除session
	this.DelSession("userName")
	//跳转登录页面
	this.Redirect("/",302)
}




