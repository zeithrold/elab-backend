package apply

import (
	"context"
	"elab-backend/service"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type SetRoomSelectionRequest struct {
	// Id 是房间的唯一标识符。
	Id string `json:"id,omitempty"`
}

type GetRoomListResponse struct {
	Rooms []RoomListItem `json:"rooms"`
}

type RoomListItem struct {
	// Id 是房间的唯一标识符。
	Id string `json:"id"`
	// Name 是房间的名称。
	Name string `json:"name"`
	// Time 是面试时间。
	Time *time.Time `json:"time"`
	// Capacity 是房间的容量。
	Capacity int `json:"capacity"`
	// Occupancy 是房间的占用情况。
	Occupancy int `json:"occupancy"`
	// Location 是房间地点。
	Location string `json:"location"`
}

type GetRoomDateListResponse struct {
	// DateList 是房间的日期列表。
	Dates []string `json:"dates"`
}

// Room 是面试房间的数据库模型。
type Room struct {
	gorm.Model
	// RoomId 是房间的唯一标识符。
	RoomId string `gorm:"type:varchar(36)"`
	// Name 是房间的名称。
	Name string `gorm:"type:varchar(255)"`
	// Time 是面试时间。
	Time *time.Time `gorm:"type:datetime"`
	// Capacity 是房间的容量。
	Capacity int `gorm:"type:int"`
	// Occupancy 是房间的占用情况。
	Occupancy int `gorm:"type:int"`
	// Location 是房间地点。
	Location string `gorm:"type:varchar(255)"`
	// Available 是房间是否可用。
	Available *bool `gorm:"type:bool"`
}

type RoomNotFoundError struct{}

func (e *RoomNotFoundError) Error() string {
	return "房间不存在"
}

type RoomFullError struct{}

func (e *RoomFullError) Error() string {
	return "房间已满"
}

type DuplicateSelectionError struct{}

func (e *DuplicateSelectionError) Error() string {
	return "重复选择房间"
}

type SelectionNotFoundError struct{}

func (e *SelectionNotFoundError) Error() string {
	return "用户未选择房间"
}

// GetRoomList 获取房间列表。
// `
// ctx 是上下文。
// date 是面试日期，格式为“YYYY-MM-DD”。
func GetRoomList(ctx context.Context, date string) *GetRoomListResponse {
	var rooms []Room
	srv := service.GetService()
	timeStart := date + " 00:00:00"
	timeEnd := date + " 23:59:59"
	slog.Debug("model.GetRoomList: 正在获取房间列表", "timeStart", timeStart, "timeEnd", timeEnd)
	err := srv.DB.WithContext(ctx).Model(&Room{}).Where(&Room{
		Available: &[]bool{true}[0],
	}).Where("time BETWEEN ? AND ?", timeStart, timeEnd).Find(&rooms).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	var res []RoomListItem
	for _, room := range rooms {
		res = append(res, RoomListItem{
			Id:        room.RoomId,
			Name:      room.Name,
			Time:      room.Time,
			Capacity:  room.Capacity,
			Occupancy: room.Occupancy,
			Location:  room.Location,
		})
	}
	return &GetRoomListResponse{
		Rooms: res,
	}
}

// GetRoomDateList 获取房间日期列表。
//
// ctx 是上下文。
func GetRoomDateList(ctx context.Context) *GetRoomDateListResponse {
	var rooms []Room
	srv := service.GetService()
	slog.Debug("model.GetRoomDateList: 正在获取房间日期列表")
	err := srv.DB.WithContext(ctx).Model(&Room{}).Where(&Room{
		Available: &[]bool{true}[0],
	}).Find(&rooms).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	var dates []string
	for _, room := range rooms {
		dates = append(dates, room.Time.Format("2006-01-02"))
	}
	// 去重
	dates = removeDuplicateElement(dates)
	return &GetRoomDateListResponse{
		Dates: dates,
	}
}

func removeDuplicateElement(a []string) []string {
	result := make([]string, 0)
	temp := map[string]struct{}{}
	for _, item := range a {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// Selection 是用户的房间选择的数据库模型。
type Selection struct {
	gorm.Model
	// OpenId 是用户的OpenId。
	OpenId string `gorm:"type:varchar(40)"`
	// RoomId 是房间的唯一标识符。
	RoomId string `gorm:"type:varchar(36)"`
}

// SetSelection 设置用户的房间选择。
//
// ctx 是上下文。
// openid 是用户的Openid。
// roomId 是房间的唯一标识符。
func SetSelection(ctx context.Context, openid string, roomId string) error {
	slog.Debug("model.SetSelection: 正在设置用户的房间选择", "openid", openid, "roomId", roomId)
	srv := service.GetService()
	// 检测房间是否存在
	isRoomExists := CheckIsRoomExists(ctx, roomId)
	if !isRoomExists {
		slog.Error("model.SetSelection: 房间不存在", "roomId", roomId)
		return &RoomNotFoundError{}
	}
	// 先获取用户是否已经选择了房间
	selectedRoomId, isAlreadySelected := CheckIsAlreadySelected(ctx, openid)
	if isAlreadySelected {
		slog.Debug("model.SetSelection: 用户已经选择了房间，可能为更改选择", "openid", openid)
		// 先确认前后房间是否相同
		if selectedRoomId == roomId {
			slog.Error("model.SetSelection: 用户选择的房间与之前相同，无需更改", "openid", openid)
			return &DuplicateSelectionError{}
		}
		slog.Debug("model.SetSelection: 用户选择的房间与之前不同，正在移除之前的选择", "openid", openid, "roomId", roomId)
		// 前后房间不同，先移除之前的选择
		err := RemoveSelection(ctx, openid, selectedRoomId)
		if err != nil {
			slog.Error("model.SetSelection: 移除用户之前的选择失败", "openid", openid, "roomId", roomId)
			panic(err)
		}
	}
	targetRoom := Room{
		RoomId:    roomId,
		Available: &[]bool{true}[0],
	}
	// 获取房间的信息
	err := srv.DB.WithContext(ctx).Model(&Room{}).Where(&targetRoom).First(&targetRoom).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	// 检测房间是否已满
	isFull := targetRoom.Occupancy >= targetRoom.Capacity
	if isFull {
		slog.Error("model.SetSelection: 房间已满", "roomId", roomId)
		return &RoomFullError{}
	}
	selection := Selection{
		OpenId: openid,
		RoomId: roomId,
	}
	err = srv.DB.WithContext(ctx).Create(&selection).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	err = srv.DB.WithContext(ctx).Model(&Room{}).Where(&Room{
		RoomId:    roomId,
		Available: &[]bool{true}[0],
	}).Updates(&Room{
		Occupancy: targetRoom.Occupancy + 1,
	}).Error
	return nil
}

func CheckIsRoomExists(ctx context.Context, roomId string) bool {
	slog.Debug("model.CheckIsRoomExists: 正在检查房间是否可用", "roomId", roomId)
	srv := service.GetService()
	targetRoom := Room{
		RoomId:    roomId,
		Available: &[]bool{true}[0],
	}
	err := srv.DB.WithContext(ctx).Model(&Room{}).Where(&targetRoom).First(&targetRoom).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("model.CheckIsRoomExists: 房间不存在", "roomId", roomId)
			return false
		} else {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		}
	}
	return true
}

func CheckIsAlreadySelected(ctx context.Context, openid string) (string, bool) {
	slog.Debug("model.CheckIsAlreadySelected: 正在检查用户是否已经选择了房间", "openid", openid)
	srv := service.GetService()
	selection := Selection{OpenId: openid}
	err := srv.DB.WithContext(ctx).Model(&Selection{}).Where(&selection).First(&selection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Debug("model.CheckIsAlreadySelected: 用户未选择房间", "openid", openid)
			return "", false
		} else {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		}
	}
	slog.Debug("model.CheckIsAlreadySelected: 用户已选择房间", "openid", openid)
	return selection.RoomId, true
}

func ClearSelection(ctx context.Context, openid string) error {
	slog.Debug("model.ClearSelection: 正在清除用户的房间选择", "openid", openid)
	srv := service.GetService()
	selection := Selection{OpenId: openid}
	slog.Debug("model.ClearSelection: 正在获取用户的房间选择", "openid", openid)
	err := srv.DB.WithContext(ctx).Model(&Selection{}).Where(&selection).First(&selection).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	roomId := selection.RoomId
	slog.Debug("model.ClearSelection: 正在移除用户的房间选择", "openid", openid)
	err = srv.DB.WithContext(ctx).Model(&Selection{}).Where(&selection).Delete(&Selection{}).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	slog.Debug("model.ClearSelection: 正在更新房间占用情况", "openid", openid)
	targetRoom := Room{
		RoomId:    roomId,
		Available: &[]bool{true}[0],
	}
	err = srv.DB.WithContext(ctx).Model(&Room{}).Where(&targetRoom).First(&targetRoom).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	targetRoom.Occupancy--
	err = srv.DB.Save(&targetRoom).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	return nil
}

func RemoveSelection(ctx context.Context, openid string, roomId string) error {
	slog.Debug("model.RemoveSelection: 正在移除用户的房间选择", "openid", openid, "roomId", roomId)
	srv := service.GetService()
	selection := Selection{
		OpenId: openid,
		RoomId: roomId,
	}
	err := srv.DB.WithContext(ctx).Where(&selection).Delete(&selection).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	slog.Debug("model.RemoveSelection 正在更新房间占用情况", "roomId", roomId)
	targetRoom := Room{
		RoomId:    roomId,
		Available: &[]bool{true}[0],
	}
	err = srv.DB.WithContext(ctx).Model(&Room{}).Where(&targetRoom).First(&targetRoom).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	slog.Debug("model.RemoveSelection 房间目前数量", "roomId", roomId, "occupancy", targetRoom.Occupancy)
	targetRoom.Occupancy--
	err = srv.DB.WithContext(ctx).Save(&targetRoom).Error
	if err != nil {
		slog.Error("调用ORM失败。", "error", err)
		panic(err)
	}
	return nil
}

func CheckIsSelectionExists(ctx context.Context, openid string) bool {
	slog.Debug("model.CheckIsSelectionExists: 正在检查用户是否已经选择了房间", "openid", openid)
	srv := service.GetService()
	selection := Selection{OpenId: openid}
	err := srv.DB.WithContext(ctx).Model(&Selection{}).Where(&selection).First(&selection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Debug("model.CheckIsSelectionExists: 用户未选择房间", "openid", openid)
			return false
		} else {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		}
	}
	slog.Debug("model.CheckIsSelectionExists: 用户已选择房间", "openid", openid)
	return true
}

func GetSelection(ctx context.Context, openid string) (*Selection, error) {
	slog.Debug("model.GetSelection: 正在获取用户的房间选择", "openid", openid)
	srv := service.GetService()
	selection := Selection{
		OpenId: openid,
	}
	err := srv.DB.WithContext(ctx).Model(&Selection{}).Where(&selection).First(&selection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Debug("model.GetSelection: 用户未选择房间", "openid", openid)
			return nil, &SelectionNotFoundError{}
		} else {
			slog.Error("调用ORM失败。", "error", err)
			panic(err)
		}
	}
	return &selection, nil
}
