package apply

import "context"

// GetStatusResponse 是获取用户状态的响应。
type GetStatusResponse struct {
	// Ticket 是用户是否已经提交申请表。
	Ticket bool `json:"ticket"`
	// RoomSelection 是用户是否已经选择房间。
	RoomSelection bool `json:"room_selection"`
	// TextForm 是用户是否已经填写文本表单。
	TextForm bool `json:"textform"`
}

// GetStatus 获取用户的状态。
//
// ctx 是上下文。
// openid 是用户的Openid。
func GetStatus(ctx context.Context, openid string) *GetStatusResponse {
	return &GetStatusResponse{
		Ticket:        CheckIsTicketExists(ctx, openid),
		RoomSelection: CheckIsSelectionExists(ctx, openid),
		TextForm:      CheckIsTextFormSubmitted(ctx, openid),
	}
}
