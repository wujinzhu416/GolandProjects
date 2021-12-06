package main

import (
	"github.com/astaxie/beego"
	_ "myproject/models"
	_ "myproject/routers"
)

func main() {
	beego.AddFuncMap("ShowPrePage",HandlePrePage)
	beego.AddFuncMap("ShowAfterPage",HandleAfterPage)
	beego.Run()
}

func HandlePrePage(data int)(int){
	pageIndex := data - 1

	return pageIndex
}

func HandleAfterPage(data int)(int){
	pageIndex := data + 1
	return pageIndex
}
