package main

import (
	"html/template"
	"net/http"

	"github.com/Codgi-123/we-wiki/app"
	"github.com/Codgi-123/we-wiki/app/controllers"
	systemControllers "github.com/Codgi-123/we-wiki/app/modules/system/controllers"
	"github.com/Codgi-123/we-wiki/app/utils"
	"github.com/astaxie/beego"
)

func init() {
	initRouter()
}

func initRouter() {
	// router
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.RouterCaseSensitive = false

	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/author", &controllers.AuthorController{}, "*:Index")
	beego.AutoRouter(&controllers.AuthorController{})
	beego.AutoRouter(&controllers.MainController{})
	beego.AutoRouter(&controllers.SpaceController{})
	beego.AutoRouter(&controllers.CollectionController{})
	beego.AutoRouter(&controllers.FollowController{})
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.DocumentController{})
	beego.AutoRouter(&controllers.PageController{})
	beego.AutoRouter(&controllers.ImageController{})

	systemNamespace := beego.NewNamespace("/system",
		beego.NSAutoRouter(&systemControllers.MainController{}),
		beego.NSAutoRouter(&systemControllers.ProfileController{}),
		beego.NSAutoRouter(&systemControllers.UserController{}),
		beego.NSAutoRouter(&systemControllers.RoleController{}),
		beego.NSAutoRouter(&systemControllers.PrivilegeController{}),
		beego.NSAutoRouter(&systemControllers.SpaceController{}),
		beego.NSAutoRouter(&systemControllers.Space_UserController{}),
	)
	beego.AddNamespace(systemNamespace)

	beego.ErrorHandler("404", http_404)
	beego.ErrorHandler("500", http_500)

	// add template func
	beego.AddFuncMap("dateFormat", utils.Date.Format)
}

func http_404(rw http.ResponseWriter, req *http.Request) {
	t, _ := template.New("404.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/error/404.html")
	data := make(map[string]interface{})
	data["content"] = "page not found"
	data["copyright"] = app.CopyRight
	t.Execute(rw, data)
}

func http_500(rw http.ResponseWriter, req *http.Request) {
	t, _ := template.New("500.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/error/500.html")
	data := make(map[string]interface{})
	data["content"] = "Server Error"
	data["copyright"] = app.CopyRight
	t.Execute(rw, data)
}
