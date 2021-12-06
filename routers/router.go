package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"myproject/controllers"
)

func init() {
	//beego.Router("/", &controllers.MainController{})

	beego.InsertFilter("/Article/*",beego.BeforeRouter,filtFunc)
	beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandlePost")

	beego.Router("/",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/logout",&controllers.UserController{},"get:Logout")

	beego.Router("/Article/showArticleList",&controllers.ArticleController{},"get,post:ShowArticleList")
	beego.Router("/Article/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/Article/showContent",&controllers.ArticleController{},"get:ShowContent")

	beego.Router("/Article/deleteArticle",&controllers.ArticleController{},"get:HandleDelete")
	beego.Router("/Article/updateArticle",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")

	//添加类型
	beego.Router("/Article/addArticleType", &controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
}

var filtFunc = func(ctx *context.Context){
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/")
	}
}