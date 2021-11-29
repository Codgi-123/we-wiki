package controllers

import (
	"fmt"
	"strings"

	"github.com/Codgi-123/we-wiki/app/models"
	"github.com/Codgi-123/we-wiki/app/utils"
	"github.com/astaxie/beego"
)

type AuthorController struct {
	BaseController
}

// login index
func (this *AuthorController) Index() {

	// is open auth login
	ssoOpen := models.ConfigModel.GetConfigValueByKey(models.ConfigKeyAuthLogin, "0")
	this.Data["sso_open"] = ssoOpen
	this.viewLayout("author/login", "author")
}

// sign up
func (this *AuthorController) SignUp() {

	if !this.IsPost() {
		this.ViewError("请求方式有误！")
	}
	username := strings.TrimSpace(this.GetString("username"))
	password := strings.TrimSpace(this.GetString("password"))
	email := strings.TrimSpace(this.GetString("email"))
	phone := strings.TrimSpace(this.GetString("phone"))

	if username == "" {
		this.jsonError("系统用户名不能为空！")
	}
	if strings.Contains(username, "_") {
		this.jsonError("系统用户名不合法！")
	}
	if password == "" {
		this.jsonError("密码不能为空！")
	}

	user, err := models.UserModel.InsertNew(username, password, phone, email)
	if err != nil {
		this.jsonError("注册出错")
	}
	_ = user

	this.jsonSuccess("注册成功！", nil, "/author/index")
}

// login
func (this *AuthorController) Login() {

	if !this.IsPost() {
		this.ViewError("请求方式有误！")
	}
	username := strings.TrimSpace(this.GetString("username"))
	password := strings.TrimSpace(this.GetString("password"))

	if username == "" {
		this.jsonError("系统用户名不能为空！")
	}
	if strings.Contains(username, "_") {
		this.jsonError("系统用户名不合法！")
	}
	if password == "" {
		this.jsonError("密码不能为空！")
	}

	user, err := models.UserModel.GetUserByUsername(username)
	if err != nil {
		this.jsonError("登录出错")
	}
	if len(user) == 0 {
		this.jsonError("用户名或密码错误!")
	}
	if user["is_forbidden"] == fmt.Sprintf("%d", models.User_Forbidden_True) {
		this.jsonError("用户已被禁用!")
	}

	password = utils.Encrypt.Md5Encode(password)

	if user["password"] != password {
		this.jsonError("用户名或密码错误!")
	}

	// save session
	this.SetSession("author", user)
	// save cookie
	identify := utils.Encrypt.Md5Encode(this.Ctx.Request.UserAgent() + password)
	passportValue := utils.Encrypt.Base64Encode(username + "@" + identify)
	passport := beego.AppConfig.String("author::passport")
	cookieExpired, _ := beego.AppConfig.Int64("author::cookie_expired")
	this.Ctx.SetCookie(passport, passportValue, cookieExpired)

	this.Ctx.Request.PostForm.Del("password")

	this.jsonSuccess("登录成功！", nil, "/main/index")
}

//logout
func (this *AuthorController) Logout() {
	passport := beego.AppConfig.String("author::passport")
	this.Ctx.SetCookie(passport, "")
	this.SetSession("author", nil)
	this.DelSession("author")

	this.Redirect("/", 302)
}
