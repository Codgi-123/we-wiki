package controllers

import (
	"strings"

	"github.com/phachon/mm-wiki/app/models"
)

type MainController struct {
	BaseController
}

func (this *MainController) Index() {

	collectDocs, err := models.CollectionModel.GetCollectionsByUserIdAndType(this.UserId, models.Collection_Type_Doc)
	if err != nil {
		this.ViewError("查找收藏文档错误！")
	}
	docIds := []string{}
	for _, collectDoc := range collectDocs {
		docIds = append(docIds, collectDoc["resource_id"])
	}

	documents, err := models.DocumentModel.GetDocumentsByDocumentIds(docIds)
	if err != nil {
		this.ViewError("查找收藏文档错误！")
	}

	this.Data["documents"] = documents
	this.Data["count"] = len(documents)
	this.viewLayout("main/index", "main")
}

func (this *MainController) Default() {

	page, _ := this.GetInt("page", 1)
	number, _ := this.GetInt("number", 10)
	maxPage := 10
	if page >= maxPage {
		page = maxPage
	}
	//number := 8
	limit := (page - 1) * number

	userId := this.UserId

	logDocuments, err := models.LogDocumentModel.GetLogDocumentsByLimit(userId, limit, number)
	if err != nil {
		this.ViewError("查找更新文档列表失败！")
	}

	count, err := models.LogDocumentModel.CountLogDocuments()
	if err != nil {
		this.ViewError("查找更新文档列表失败！")
	}
	if count >= int64(maxPage*number) {
		count = int64(maxPage * number)
	}

	userIds := []string{}
	docIds := []string{}
	for _, logDocument := range logDocuments {
		userIds = append(userIds, logDocument["user_id"])
		docIds = append(docIds, logDocument["document_id"])
	}
	users, err := models.UserModel.GetUsersByUserIds(userIds)
	if err != nil {
		this.ViewError("查找更新文档列表失败！")
	}
	docs, err := models.DocumentModel.GetAllDocumentsByDocumentIds(docIds)
	if err != nil {
		this.ViewError("查找更新文档列表失败！")
	}
	for _, logDocument := range logDocuments {
		logDocument["username"] = ""
		for _, user := range users {
			if logDocument["user_id"] == user["user_id"] {
				logDocument["username"] = user["username"]
				break
			}
		}
		for _, doc := range docs {
			if logDocument["document_id"] == doc["document_id"] {
				logDocument["document_name"] = doc["name"]
				logDocument["document_type"] = doc["type"]
				break
			}
		}
	}

	// main title config
	mainTitle := models.ConfigModel.GetConfigValueByKey(models.ConfigKeyMainTitle, "")
	mainDescription := models.ConfigModel.GetConfigValueByKey(models.ConfigKeyMainDescription, "")

	this.Data["panel_title"] = mainTitle
	this.Data["panel_description"] = mainDescription
	this.Data["logDocuments"] = logDocuments
	this.SetPaginator(number, count)
	this.viewLayout("main/default", "default")
}

// 搜索，支持根据标题和内容搜索
func (this *MainController) Search() {

	keyword := strings.TrimSpace(this.GetString("keyword", ""))
	searchType := this.GetString("search_type", "content")

	this.Data["search_type"] = searchType
	this.Data["keyword"] = keyword
	this.Data["count"] = 0
	if keyword == "" {
		this.viewLayout("main/search", "default")
		return
	}
	var documents = []map[string]string{}
	var err error
	// 获取该用户有权限的空间
	publicSpaces, err := models.SpaceModel.GetSpacesByVisitLevel(models.Space_VisitLevel_Public)
	if err != nil {
		this.ViewError("搜索文档错误！")
	}
	spaceUsers, err := models.SpaceUserModel.GetSpaceUsersByUserId(this.UserId)
	if err != nil {
		this.ViewError("搜索文档错误！")
	}
	spaceIdsMap := make(map[string]bool)
	for _, publicSpace := range publicSpaces {
		spaceIdsMap[publicSpace["space_id"]] = true
	}
	for _, spaceUser := range spaceUsers {
		if _, ok := spaceIdsMap[spaceUser["space_id"]]; !ok {
			spaceIdsMap[spaceUser["space_id"]] = true
		}
	}
	searchDocContents := make(map[string]string)
	searchType = "title"
	documents, err = models.DocumentModel.GetDocumentsByLikeName(keyword)
	if err != nil {
		this.ViewError("搜索文档错误！")
	}
	// 过滤一下没权限的空间
	realDocuments := []map[string]string{}
	for _, document := range documents {
		spaceId, _ := document["space_id"]
		documentId, _ := document["document_id"]
		if _, ok := spaceIdsMap[spaceId]; !ok {
			continue
		}
		if searchType != "title" {
			searchContent, ok := searchDocContents[documentId]
			if !ok || searchContent == "" {
				continue
			}
			document["search_content"] = searchContent
		}
		realDocuments = append(realDocuments, document)
	}

	this.Data["search_type"] = searchType
	this.Data["keyword"] = keyword
	this.Data["documents"] = realDocuments
	this.Data["count"] = len(realDocuments)
	this.viewLayout("main/search", "default")
}
