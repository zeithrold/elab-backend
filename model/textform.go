package model

import (
	"context"
	"elab-backend/service"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log/slog"
)

// TextForm 是用户的文字表单。
type TextForm struct {
	gorm.Model
	// OpenId 是用户的OpenId。
	OpenId string `gorm:"type:varchar(36)"`
	// Question 是用户需要回答的问题。
	QuestionId string `gorm:"type:varchar(1024)"`
	// Answer 是用户的回答。
	Answer string `gorm:"type:varchar(1024)"`
}

// Question 是用户需要回答的问题列表
type Question struct {
	gorm.Model
	// Id 是用户需要回答的问题ID。
	Id string `gorm:"type:varchar(36)"`
	// Question 是问题标题。
	Question string `gorm:"type:varchar(1024)"`
	// Text 是问题的文字描述。
	Text string `gorm:"type:varchar(1024)"`
}

type UpdateTextFormRequest struct {
	TextForms []struct {
		Id     string `json:"id"`
		Answer string `json:"answer"`
	} `json:"text_forms"`
}

// GetQuestionList 获取问题列表。
//
// ctx 是上下文。
func GetQuestionList(ctx context.Context) []Question {
	slog.Debug("model.GetQuestionList: 正在获取问题列表")
	srv := service.GetService()
	var questions []Question
	err := srv.DB.WithContext(ctx).Model(&questions).Find(&questions).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	return questions
}

// UpdateTextForm 更新用户的文本表单。
//
// ctx 是上下文。
// openid 是用户的Openid。
// request 是用户的请求。
func UpdateTextForm(ctx context.Context, openid string, request UpdateTextFormRequest) {
	slog.Debug("model.UpdateTextForm: 正在更新文本表单", "openid", openid, "len", len(request.TextForms))
	srv := service.GetService()
	slog.Debug("model.UpdateTextForm: 正在检查用户是否已经填写了文本表单", "openid", openid)
	isTextFormExists := CheckIsTextFormExists(ctx, openid)
	if isTextFormExists {
		slog.Debug("model.UpdateTextForm: 用户已经填写了文本表单，正在清除")
		ClearQuestion(ctx, openid)
	}
	for _, v := range request.TextForms {
		slog.Debug("model.UpdateTextForm: 正在更新文本表单", "openid", openid, "questionId", v.Id, "answer", v.Answer)
		err := srv.DB.WithContext(ctx).Model(&TextForm{}).Create(&TextForm{
			OpenId:     openid,
			QuestionId: v.Id,
			Answer:     v.Answer,
		}).Error
		if err != nil {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		}
	}
	slog.Debug("model.UpdateTextForm: 更新文本表单成功", "openid", openid)
}

// CheckIsTextFormExists 检查用户是否已经填写了文本表单。
//
// ctx 是上下文。
// openid 是用户的Openid。
func CheckIsTextFormExists(ctx context.Context, openid string) bool {
	slog.Debug("model.CheckIsTextFormExists: 正在检查用户是否已经填写了文本表单", "openid", openid)
	srv := service.GetService()
	var textForm TextForm
	err := srv.DB.WithContext(ctx).Model(&textForm).Where(&TextForm{
		OpenId: openid,
	}).First(&textForm).Error
	if err != nil {
		isNotExist := errors.Is(err, gorm.ErrRecordNotFound)
		if !isNotExist {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		}
		slog.Debug("model.CheckIsTextFormExists: 用户未填写文本表单", "openid", openid)
		return false
	}
	slog.Debug("model.CheckIsTextFormExists: 用户已填写文本表单", "openid", openid)
	return true
}

// ClearQuestion 清除用户的问题。
//
// ctx 是上下文。
// openid 是用户的Openid。
func ClearQuestion(ctx context.Context, openid string) {
	slog.Debug("model.ClearQuestion: 正在清除用户的文本表单", "openid", openid)
	srv := service.GetService()
	err := srv.DB.WithContext(ctx).Model(&TextForm{}).Where(&TextForm{
		OpenId: openid,
	}).Delete(&TextForm{}).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
}
