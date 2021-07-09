package controllers

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/phachon/mm-wiki/app/models"
)

type UserController struct {
	BaseController
}

func (this *UserController) Add() {

	roles := []map[string]string{}
	var err error

	if this.IsRoot() {
		roles, err = models.RoleModel.GetRoles()
	} else {
		roles, err = models.RoleModel.GetRolesNotContainRoot()
	}
	if err != nil {
		this.ViewError("获取用户角色失败！")
	}
	this.Data["roles"] = roles
	this.viewLayout("user/form", "user")
}

func (this *UserController) Save() {

	if !this.IsPost() {
		this.ViewError("请求方式有误！", "/system/user/add")
	}
	username := strings.TrimSpace(this.GetString("username", ""))
	password := strings.TrimSpace(this.GetString("password", ""))
	email := strings.TrimSpace(this.GetString("email", ""))
	phone := strings.TrimSpace(this.GetString("phone", ""))
	roleId := strings.TrimSpace(this.GetString("role_id", ""))
	this.Ctx.Request.PostForm.Del("password")

	v := validation.Validation{}
	if username == "" {
		this.jsonError("用户名不能为空！")
	}
	if !v.AlphaNumeric(username, "username").Ok {
		this.jsonError("用户名格式不正确！")
	}
	if password == "" {
		this.jsonError("密码不能为空！")
	}
	if email == "" {
		this.jsonError("邮箱不能为空！")
	}
	if !v.Email(email, "email").Ok {
		this.jsonError("邮箱格式不正确！")
	}
	if phone == "" {
		this.jsonError("手机号不能为空！")
	}
	if roleId == "" {
		this.jsonError("没有选择角色！")
	}

	ok, err := models.UserModel.CreateNew(username, password, phone, email, roleId)
	if err != nil {
		this.jsonError("添加用户失败！")
	}
	_ = ok

	if err != nil {
		this.jsonError("添加用户失败")
	}

	this.jsonSuccess("添加用户成功", nil, "/system/user/list")
}

func (this *UserController) List() {

	keywords := map[string]string{}
	page, _ := this.GetInt("page", 1)
	username := strings.TrimSpace(this.GetString("username", ""))
	roleId := strings.TrimSpace(this.GetString("role_id", ""))
	number, _ := this.GetRangeInt("number", 20, 10, 100)

	if username != "" {
		keywords["username"] = username
	}
	if roleId != "" {
		keywords["role_id"] = roleId
	}

	limit := (page - 1) * number
	var err error
	var count int64
	var users []map[string]string
	if len(keywords) != 0 {
		count, err = models.UserModel.CountUsersByKeywords(keywords)
		users, err = models.UserModel.GetUsersByKeywordsAndLimit(keywords, limit, number)
	} else {
		count, err = models.UserModel.CountUsers()
		users, err = models.UserModel.GetUsersByLimit(limit, number)
	}
	if err != nil {
		this.ViewError("获取用户列表失败", "/system/main/index")
	}

	var roleIds = []string{}
	if roleId != "" {
		roleIds = append(roleIds, roleId)
	} else {
		for _, user := range users {
			roleIds = append(roleIds, user["role_id"])
		}
	}
	roles, err := models.RoleModel.GetRoleByRoleIds(roleIds)
	if err != nil {
		this.ViewError("获取用户列表失败!", "/system/main/index")
	}
	var roleUsers = []map[string]string{}
	for _, user := range users {
		roleUser := user
		for _, role := range roles {
			if role["role_id"] == user["role_id"] {
				roleUser["role_name"] = role["name"]
				break
			}
		}
		roleUsers = append(roleUsers, roleUser)
	}

	allRoles, err := models.RoleModel.GetRoles()
	if err != nil {
		this.ViewError("获取用户列表失败！", "/system/main/index")
	}
	this.Data["users"] = roleUsers
	this.Data["username"] = username
	this.Data["roleId"] = roleId
	this.Data["roles"] = allRoles
	this.SetPaginator(number, count)
	this.viewLayout("user/list", "user")
}

func (this *UserController) Edit() {

	userId := this.GetString("user_id", "")
	if userId == "" {
		this.ViewError("用户不存在！", "/system/user/list")
	}

	user, err := models.UserModel.GetUserByUserId(userId)
	if err != nil {
		this.ViewError("查找用户出错！", "/system/user/list")
	}
	if len(user) == 0 {
		this.ViewError("用户不存在！", "/system/user/list")
	}
	// 登录非 root 用户不能修改 root 用户信息
	if user["role_id"] == fmt.Sprintf("%d", models.Role_Root_Id) && !this.IsRoot() {
		this.ViewError("没有权限修改！", "/system/user/list")
	}

	roles := []map[string]string{}
	if this.IsRoot() {
		roles, err = models.RoleModel.GetRoles()
	} else {
		roles, err = models.RoleModel.GetRolesNotContainRoot()
	}
	if err != nil {
		this.ViewError("获取用户角色失败！")
	}

	this.Data["user"] = user
	this.Data["roles"] = roles
	this.viewLayout("user/edit", "user")
}

func (this *UserController) Modify() {

	if !this.IsPost() {
		this.ViewError("请求方式有误！", "/system/user/list")
	}
	userId := strings.TrimSpace(this.GetString("user_id", ""))
	email := strings.TrimSpace(this.GetString("email", ""))
	phone := strings.TrimSpace(this.GetString("phone", ""))
	roleId := strings.TrimSpace(this.GetString("role_id", ""))
	password := strings.TrimSpace(this.GetString("password", ""))
	this.Ctx.Request.PostForm.Del("password")

	v := validation.Validation{}
	if email == "" {
		this.jsonError("邮箱不能为空！")
	}
	if !v.Email(email, "email").Ok {
		this.jsonError("邮箱格式不正确！")
	}
	if phone == "" {
		this.jsonError("手机号不能为空！")
	}

	user, err := models.UserModel.GetUserByUserId(userId)
	if err != nil {
		this.jsonError("修改用户出错！")
	}
	if len(user) == 0 {
		this.jsonError("用户不存在！")
	}
	if user["role_id"] == fmt.Sprintf("%d", models.Role_Root_Id) {
		roleId = fmt.Sprintf("%d", models.Role_Root_Id)
	}
	if password == "" {
		this.jsonError("密码不能为空！")
	}
	if user["role_id"] == fmt.Sprintf("%d", models.Role_Root_Id) && !this.IsRoot() {
		this.jsonError("没有权限修改！")
	}

	updateUser := map[string]interface{}{
		"email":   email,
		"phone":   phone,
		"role_id": roleId,
	}
	if password != "" {
		updateUser["password"] = models.UserModel.EncodePassword(password)
	}
	if roleId != "" {
		updateUser["role_id"] = roleId
	}
	_, err = models.UserModel.Update(userId, updateUser)
	if err != nil {
		this.jsonError("修改用户失败")
	}
	this.jsonSuccess("修改用户成功", nil, "/system/user/list")
}

func (this *UserController) Forbidden() {

	if !this.IsPost() {
		this.ViewError("请求方式有误！", "/system/user/list")
	}
	userId := this.GetString("user_id", "")
	if userId == "" {
		this.jsonError("用户不存在")
	}

	user, err := models.UserModel.GetUserByUserId(userId)
	if err != nil {
		this.jsonError("屏蔽用户失败")
	}
	if len(user) == 0 {
		this.jsonError("用户不存在")
	}
	if user["role_id"] == fmt.Sprintf("%d", models.Role_Root_Id) {
		this.jsonError("不能操作超级管理员")
	}
	_, err = models.UserModel.Update(userId, map[string]interface{}{
		"is_forbidden": models.User_Forbidden_True,
	})
	if err != nil {
		this.jsonError("屏蔽用户失败")
	}

	this.jsonSuccess("屏蔽用户成功", nil, "/system/user/list")
}

func (this *UserController) Recover() {

	if !this.IsPost() {
		this.ViewError("请求方式有误！", "/system/user/list")
	}
	userId := this.GetString("user_id", "")
	if userId == "" {
		this.jsonError("用户不存在")
	}

	user, err := models.UserModel.GetUserByUserId(userId)
	if err != nil {
		this.jsonError("恢复用户失败")
	}
	if len(user) == 0 {
		this.jsonError("用户不存在")
	}
	if user["role_id"] == fmt.Sprintf("%d", models.Role_Root_Id) {
		this.jsonError("不能操作超级管理员")
	}
	_, err = models.UserModel.Update(userId, map[string]interface{}{
		"is_forbidden": models.User_Is_Forbidden_False,
	})
	if err != nil {
		this.jsonError("恢复用户失败")
	}

	this.jsonSuccess("恢复用户成功", nil, "/system/user/list")
}

func (this *UserController) Info() {

	userId := this.GetString("user_id", "")
	if userId == "" {
		this.ViewError("用户不存在！", "/system/user/list")
	}

	user, err := models.UserModel.GetUserByUserId(userId)
	if err != nil {
		this.ViewError("查找用户出错！", "/system/user/list")
	}
	if len(user) == 0 {
		this.ViewError("用户不存在！", "/system/user/list")
	}
	role, err := models.RoleModel.GetRoleByRoleId(user["role_id"])
	if err != nil {
		this.ViewError("查找用户出错！", "/system/user/list")
	}
	this.Data["user"] = user
	this.Data["role"] = role
	this.viewLayout("user/info", "user")
}
