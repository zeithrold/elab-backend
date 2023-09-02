package apply

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
	OpenId string `gorm:"type:varchar(40)"`
	// Question 是用户需要回答的问题。
	QuestionId string `gorm:"type:varchar(1024)"`
	// Answer 是用户的回答。
	Answer string `gorm:"type:varchar(1024)"`
	// Submitted 是用户是否已经提交申请表。
	Submitted *bool `gorm:"type:bool"`
}

// Question 是用户需要回答的问题列表
type Question struct {
	gorm.Model
	// QuestionId 是用户需要回答的问题ID。
	QuestionId string `gorm:"type:varchar(36)"`
	// Question 是问题标题。
	Question string `gorm:"type:varchar(1024)"`
	// Text 是问题的文字描述。
	Text string `gorm:"type:varchar(1024)"`
}

// GetQuestionListResponse 获取用户的文字表单列表。
type GetQuestionListResponse struct {
	// QuestionList 是用户需要回答的问题列表。
	Questions []QuestionListItem `json:"questions"`
}

type GetQuestionRequestUri struct {
	// Id 是用户需要回答的问题ID。
	Id string `uri:"id" binding:"required"`
}

// QuestionListItem 是用户需要回答的问题列表项。
type QuestionListItem struct {
	// Id 是用户需要回答的问题ID。
	Id string `json:"id"`
	// Question 是问题标题。
	Question string `json:"question"`
	// Text 是问题的文字描述。
	Text string `json:"text"`
	// Submitted 是用户是否已经提交申请表。
	Submitted bool `json:"submitted"`
}

type UpdateTextFormRequestUri struct {
	// Id 是用户需要回答的问题ID。
	Id string `uri:"id" binding:"required"`
}

type UpdateTextFormRequest struct {
	// Id 是用户需要回答的问题ID。
	Id string `json:"id"`
	// Answer 是用户的回答。
	Answer string `json:"answer"`
}

// GetTextFormListResponse 获取用户的文字表单列表。
type GetTextFormListResponse struct {
	TextForms []TextFormListItem `json:"textform"`
}

// TextFormListItem 是用户的文字表单列表项。
type TextFormListItem struct {
	// Id 是用户需要回答的问题ID。
	Id string `json:"id"`
	// Answer 是用户的回答。
	Answer string `json:"answer"`
	// Submitted 是用户是否已经提交申请表。
	Submitted bool `json:"submitted"`
}

// GetQuestionList 获取问题列表。
//
// ctx 是上下文。
func GetQuestionList(ctx context.Context, openid string) *GetQuestionListResponse {
	slog.Debug("model.GetQuestionList: 正在获取问题列表")
	srv := service.GetService()
	var questions []Question
	err := srv.DB.WithContext(ctx).Model(&Question{}).Find(&questions).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	textFormList := GetTextForm(ctx, openid)
	var result GetQuestionListResponse
	for _, v := range questions {
		var submitted bool
		for _, vv := range textFormList.TextForms {
			if v.QuestionId == vv.Id {
				submitted = vv.Submitted
			}
		}
		result.Questions = append(result.Questions, QuestionListItem{
			Id:        v.QuestionId,
			Question:  v.Question,
			Text:      v.Text,
			Submitted: submitted,
		})
	}
	return &result
}

func GetQuestion(ctx context.Context, openid string, questionId string) *QuestionListItem {
	slog.Debug("model.GetQuestion: 正在获取问题")
	srv := service.GetService()
	var question Question
	err := srv.DB.WithContext(ctx).Model(&Question{}).Where(&Question{QuestionId: questionId}).Find(&question).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	textFormList := GetTextForm(ctx, openid)
	var submitted bool
	for _, vv := range textFormList.TextForms {
		if question.QuestionId == vv.Id {
			submitted = vv.Submitted
		}
	}
	return &QuestionListItem{
		Id:        question.QuestionId,
		Question:  question.Question,
		Text:      question.Text,
		Submitted: submitted,
	}
}

// GetTextForm 获取用户的文本表单。
//
// ctx 是上下文。
// openid 是用户的Openid。
func GetTextForm(ctx context.Context, openid string) *GetTextFormListResponse {
	slog.Debug("model.GetTextForm: 正在获取文本表单", "openid", openid)
	isTextFormExists := CheckIsTextFormExists(ctx, openid)
	if !isTextFormExists {
		slog.Debug("model.GetTextForm: 文本表单不存在，正在初始化", "openid", openid)
		InitTextForm(ctx, openid)
	}
	srv := service.GetService()
	var textForms []TextForm
	err := srv.DB.WithContext(ctx).Model(&TextForm{}).Where(&TextForm{OpenId: openid}).Find(&textForms).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	var result GetTextFormListResponse
	for _, v := range textForms {
		result.TextForms = append(
			result.TextForms,
			TextFormListItem{
				Id:        v.QuestionId,
				Answer:    v.Answer,
				Submitted: *v.Submitted,
			})
	}
	return &result
}

func InitTextForm(ctx context.Context, openid string) {
	slog.Debug("model.InitTextForm: 正在初始化文本表单", "openid", openid)
	srv := service.GetService()
	var questions []Question
	err := srv.DB.WithContext(ctx).Model(&Question{}).Find(&questions).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	for _, v := range questions {
		slog.Debug("model.InitTextForm: 正在初始化文本表单", "openid", openid, "questionId", v.QuestionId)
		err := srv.DB.WithContext(ctx).Model(&TextForm{}).Create(
			&TextForm{OpenId: openid, QuestionId: v.QuestionId, Submitted: &[]bool{false}[0]}).Error
		if err != nil {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		}
	}
}

// UpdateTextForm 更新用户的文本表单。
//
// ctx 是上下文。
// openid 是用户的Openid。
// request 是用户的请求。
func UpdateTextForm(ctx context.Context, openid string, request *UpdateTextFormRequest) {
	slog.Debug("model.UpdateTextForm: 正在更新文本表单", "openid", openid, "questionId", request.Id)
	srv := service.GetService()
	err := srv.DB.WithContext(ctx).Model(&TextForm{}).Where(&TextForm{
		OpenId:     openid,
		QuestionId: request.Id,
	}).Updates(&TextForm{
		Answer:    request.Answer,
		Submitted: &[]bool{true}[0],
	}).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	slog.Debug("model.UpdateTextForm: 更新文本表单成功", "openid", openid)
}

// CheckIsTextFormSubmitted 检查用户是否已经填写了文本表单。
//
// ctx 是上下文。
// openid 是用户的Openid。
func CheckIsTextFormSubmitted(ctx context.Context, openid string) bool {
	slog.Debug("model.CheckIsTextFormSubmitted: 正在检查用户是否已经填写了文本表单", "openid", openid)
	srv := service.GetService()
	var textForm TextForm
	err := srv.DB.WithContext(ctx).Model(&TextForm{}).Where(&TextForm{
		OpenId: openid,
	}).First(&textForm).Error
	if err != nil {
		isNotExist := errors.Is(err, gorm.ErrRecordNotFound)
		if !isNotExist {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		}
		slog.Debug("model.CheckIsTextFormSubmitted: 没有openid为", "openid", openid, "的文本表单")
		return false
	}
	var counts int64
	err = srv.DB.WithContext(ctx).Model(&TextForm{}).Where(&TextForm{
		OpenId:    openid,
		Submitted: &[]bool{false}[0],
	}).Count(&counts).Error
	if err != nil {
		isNotExist := errors.Is(err, gorm.ErrRecordNotFound)
		if !isNotExist {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		}
	}
	if counts == 0 {
		slog.Debug("model.CheckIsTextFormSubmitted: 用户已经填写完全文本表单", "openid", openid)
		return true
	}
	slog.Debug("model.CheckIsTextFormSubmitted: 用户未填写完全文本表单", "openid", openid)
	return false
}

func CheckIsTextFormExists(ctx context.Context, openid string) bool {
	slog.Debug("model.CheckIsTextFormExists: 正在检查用户是否已经填写了文本表单", "openid", openid)
	srv := service.GetService()
	var textForm TextForm
	err := srv.DB.WithContext(ctx).Model(&TextForm{}).Where(&TextForm{
		OpenId: openid,
	}).First(&textForm).Error
	if err != nil {
		isNotExist := errors.Is(err, gorm.ErrRecordNotFound)
		if !isNotExist {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		}
		slog.Debug("model.CheckIsTextFormExists: 没有openid为", "openid", openid, "的文本表单")
		return false
	}
	slog.Debug("model.CheckIsTextFormExists: openid为", "openid", openid, "的文本表单存在")
	return true
}
