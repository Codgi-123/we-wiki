package controllers

import (
	"encoding/json"
)

type UploadResponse struct {
	Success int    `json:"success"`
	Message string `json:"message"`
	Url     string `json:"url"`
}

type ImageController struct {
	BaseController
}

func (this *ImageController) jsonError(message string) {

	uploadRes := UploadResponse{
		Success: 0,
		Message: message,
		Url:     "",
	}

	j, err := json.Marshal(uploadRes)
	if err != nil {
		this.Abort(err.Error())
	} else {
		this.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
		this.Abort(string(j))
	}
}

func (this *ImageController) jsonSuccess(message string, url string) {

	uploadRes := UploadResponse{
		Success: 1,
		Message: message,
		Url:     url,
	}

	j, err := json.Marshal(uploadRes)
	if err != nil {
		this.Abort(err.Error())
	} else {
		this.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
		this.Abort(string(j))
	}
}
