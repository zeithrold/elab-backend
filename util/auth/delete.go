package auth

import (
	"context"
	"elab-backend/model/apply"
	"elab-backend/service"
	"log/slog"
)

func DeleteAccount(ctx context.Context, openid string) error {
	svc := service.GetService()
	err := svc.DB.WithContext(ctx).Model(&apply.Ticket{}).Where(&apply.Ticket{OpenId: openid}).Delete(&apply.Ticket{}).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	exists := apply.CheckIsSelectionExists(ctx, openid)
	if exists {
		err = apply.ClearSelection(ctx, openid)
		if err != nil {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		}
	}
	err = svc.DB.WithContext(ctx).Model(&apply.TextForm{}).Where(&apply.TextForm{OpenId: openid}).Delete(&apply.TextForm{}).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	err = svc.AuthAPI.User.Delete(ctx, openid)
	if err != nil {
		slog.Error("调用Auth0 API失败。", "error", err)
		panic(err)
	}
	return nil
}
