package controllers

import (
	"github.com/Codgi-123/we-wiki/app/models"
	"github.com/astaxie/beego"
)

type PrivilegeController struct {
	BaseController
}

func (this *PrivilegeController) Add() {

	menus, _, err := models.PrivilegeModel.GetTypePrivileges()
	if err != nil {
		this.ViewError("查找权限菜单失败！")
	}

	this.Data["menus"] = menus
	this.Data["mode"] = beego.BConfig.RunMode
	this.viewLayout("privilege/form", "privilege")
}

func (this *PrivilegeController) List() {

	menus, controllers, err := models.PrivilegeModel.GetTypePrivileges()
	if err != nil {
		this.ViewError("查找权限失败！")
	}

	this.Data["menus"] = menus
	this.Data["controllers"] = controllers
	this.Data["mode"] = beego.BConfig.RunMode
	this.viewLayout("privilege/list", "privilege")
}

func (this *PrivilegeController) Edit() {

	privilegeId := this.GetString("privilege_id", "")
	if privilegeId == "" {
		this.ViewError("没有选择权限！", "/system/privilege/list")
	}

	privilege, err := models.PrivilegeModel.GetPrivilegeByPrivilegeId(privilegeId)
	if err != nil {
		this.ViewError("查找权限失败！", "/system/privilege/list")
	}
	if len(privilege) == 0 {
		this.ViewError("权限不存在！", "/system/privilege/list")
	}

	menus, _, err := models.PrivilegeModel.GetTypePrivileges()
	if err != nil {
		this.ViewError("查找权限失败！", "/system/privilege/list")
	}

	this.Data["menus"] = menus
	this.Data["privilege"] = privilege
	this.Data["mode"] = beego.BConfig.RunMode
	this.viewLayout("privilege/form", "privilege")
}
