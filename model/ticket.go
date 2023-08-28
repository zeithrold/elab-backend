package model

import (
	"context"
	"elab-backend/service"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log/slog"
)

type TicketUpdateRequest struct {
	// Name 是用户的姓名。
	Name string `json:"name"`
	// StudentId 是用户的学号。
	StudentId string `json:"student_id"`
	// ClassName 是用户的班级，以此来替代所属学院
	ClassName string `json:"class_name"`
	// Group 是用户的组别，如“软件组”、“硬件组”等。
	Group string `json:"group"`
	// Contact 是用户的联系方式，如手机号等。
	Contact string `json:"contact"`
}

// Ticket 是科中成员的申请表，用于装填基本信息。
type Ticket struct {
	gorm.Model
	// Openid 是用户的Openid。
	OpenId string `gorm:"type:varchar(36)"`
	// Name 是用户的姓名。
	Name string `gorm:"type:varchar(36)"`
	// StudentId 是用户的学号。
	StudentId string `gorm:"type:varchar(16)"`
	// ClassName 是用户的班级，以此来替代所属学院
	ClassName string `gorm:"type:varchar(16)"`
	// Group 是用户的组别，如“软件组”、“硬件组”等。
	Group string `gorm:"type:varchar(16)"`
	// Contact 是用户的联系方式，如手机号等。
	Contact string `gorm:"type:varchar(16)"`
	// Submitted 是用户是否已经提交申请表。
	Submitted bool `gorm:"type:bool"`
}

// GetTicket 获取用户的申请表。
//
// ctx 是上下文。
// openid 是用户的Openid。
func GetTicket(ctx context.Context, openid string) *Ticket {
	slog.Debug("model.GetTicket: 正在获取申请表", "openid", openid)
	srv := service.GetService()
	var ticket Ticket
	err := srv.DB.WithContext(ctx).Model(&ticket).Where(&Ticket{
		OpenId: openid,
	}).First(&ticket).Error
	if err != nil {
		isNotExist := errors.Is(err, gorm.ErrRecordNotFound)
		if !isNotExist {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		} else {
			slog.Debug("model.GetTicket: 数据库不存在用户的申请表，正在创建一个。", "openid", openid)
			InitTicket(ctx, openid)
		}
	}
	return &ticket
}

// InitTicket 初始化用户的申请表。
//
// openid 是用户的Openid。
func InitTicket(ctx context.Context, openid string) {
	slog.Debug("model.InitTicket: 正在初始化申请表", "openid", openid)
	srv := service.GetService()
	ticket := Ticket{
		OpenId:    openid,
		Submitted: false,
	}
	err := srv.DB.WithContext(ctx).Create(&ticket).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
}

// UpdateTicket 更新用户的申请表。
//
// openid 是用户的Openid。
func UpdateTicket(ctx context.Context, openid string, body *TicketUpdateRequest) {
	slog.Debug("model.UpdateTicket: 正在更新申请表", "openid", openid)
	srv := service.GetService()
	ticket := Ticket{
		OpenId:    openid,
		Name:      body.Name,
		StudentId: body.StudentId,
		ClassName: body.ClassName,
		Group:     body.Group,
		Contact:   body.Contact,
		Submitted: true,
	}
	err := srv.DB.WithContext(ctx).Model(&Ticket{}).Where(&Ticket{
		OpenId: openid,
	}).Updates(&ticket).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
}
