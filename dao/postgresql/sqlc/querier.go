// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CountAccountsByUserID(ctx context.Context, userID int64) (int32, error)
	CreateAccount(ctx context.Context, arg *CreateAccountParams) error
	CreateApplication(ctx context.Context, arg *CreateApplicationParams) error
	CreateFile(ctx context.Context, arg *CreateFileParams) (*File, error)
	CreateFriendRelation(ctx context.Context, arg *CreateFriendRelationParams) (int64, error)
	CreateGroupNotify(ctx context.Context, arg *CreateGroupNotifyParams) (*CreateGroupNotifyRow, error)
	CreateGroupRelation(ctx context.Context, arg *CreateGroupRelationParams) (int64, error)
	CreateManySetting(ctx context.Context, arg []*CreateManySettingParams) (int64, error)
	CreateMessage(ctx context.Context, arg *CreateMessageParams) (*CreateMessageRow, error)
	CreateSetting(ctx context.Context, arg *CreateSettingParams) error
	CreateUser(ctx context.Context, arg *CreateUserParams) (*User, error)
	DeleteAccount(ctx context.Context, id int64) error
	DeleteAccountsByUserID(ctx context.Context, userID int64) ([]int64, error)
	DeleteApplication(ctx context.Context, arg *DeleteApplicationParams) error
	DeleteFileByID(ctx context.Context, id int64) error
	DeleteFriendRelationsByAccountID(ctx context.Context, account1ID int64) ([]int64, error)
	DeleteGroup(ctx context.Context, relationID int64) error
	DeleteGroupNotify(ctx context.Context, id int64) error
	DeleteRelation(ctx context.Context, id int64) error
	DeleteSetting(ctx context.Context, arg *DeleteSettingParams) error
	DeleteSettingsByAccountID(ctx context.Context, accountID int64) ([]int64, error)
	DeleteUser(ctx context.Context, id int64) error
	ExistEmail(ctx context.Context, email string) (bool, error)
	ExistsAccountByID(ctx context.Context, id int64) (bool, error)
	ExistsAccountByNameAndUserID(ctx context.Context, arg *ExistsAccountByNameAndUserIDParams) (bool, error)
	ExistsApplicationByIDWithLock(ctx context.Context, arg *ExistsApplicationByIDWithLockParams) (bool, error)
	ExistsFriendRelation(ctx context.Context, arg *ExistsFriendRelationParams) (bool, error)
	ExistsFriendSetting(ctx context.Context, arg *ExistsFriendSettingParams) (bool, error)
	ExistsGroupLeaderByAccountIDWithLock(ctx context.Context, accountID int64) (bool, error)
	ExistsIsLeader(ctx context.Context, arg *ExistsIsLeaderParams) (bool, error)
	ExistsSetting(ctx context.Context, arg *ExistsSettingParams) (bool, error)
	ExistsUserByID(ctx context.Context, id int64) (bool, error)
	GetAccountByID(ctx context.Context, arg *GetAccountByIDParams) (*GetAccountByIDRow, error)
	GetAccountIDsByRelationID(ctx context.Context, relationID int64) ([]int64, error)
	GetAccountsByName(ctx context.Context, arg *GetAccountsByNameParams) ([]*GetAccountsByNameRow, error)
	GetAccountsByUserID(ctx context.Context, userID int64) ([]*GetAccountsByUserIDRow, error)
	GetAcountIDsByUserID(ctx context.Context, userID int64) ([]int64, error)
	GetAllEmail(ctx context.Context) ([]string, error)
	GetAllGroupRelation(ctx context.Context) ([]int64, error)
	GetAllRelationIDs(ctx context.Context) ([]int64, error)
	GetAllRelationOnRelation(ctx context.Context) ([]*Relation, error)
	GetApplicationByID(ctx context.Context, arg *GetApplicationByIDParams) (*Application, error)
	GetApplications(ctx context.Context, arg *GetApplicationsParams) ([]*GetApplicationsRow, error)
	GetFileByRelationID(ctx context.Context, relationID sql.NullInt64) ([]*File, error)
	GetFileByRelationIDIsNULL(ctx context.Context) ([]*GetFileByRelationIDIsNULLRow, error)
	GetFileDetailsByID(ctx context.Context, id int64) (*File, error)
	GetFileKeyByID(ctx context.Context, id int64) (string, error)
	GetFriendPinSettingsOrderByPinTime(ctx context.Context, accountID int64) ([]*GetFriendPinSettingsOrderByPinTimeRow, error)
	GetFriendRelationByID(ctx context.Context, id int64) (*GetFriendRelationByIDRow, error)
	GetFriendSettingsByName(ctx context.Context, arg *GetFriendSettingsByNameParams) ([]*GetFriendSettingsByNameRow, error)
	GetFriendSettingsOrderByName(ctx context.Context, accountID int64) ([]*GetFriendSettingsOrderByNameRow, error)
	GetFriendShowSettingsOrderByShowTime(ctx context.Context, accountID int64) ([]*GetFriendShowSettingsOrderByShowTimeRow, error)
	GetGroupAvatar(ctx context.Context, relationID sql.NullInt64) (*File, error)
	GetGroupList(ctx context.Context, accountID int64) ([]*GetGroupListRow, error)
	GetGroupMembers(ctx context.Context, relationID int64) ([]int64, error)
	GetGroupMembersByID(ctx context.Context, arg *GetGroupMembersByIDParams) ([]*GetGroupMembersByIDRow, error)
	GetGroupNotifyByID(ctx context.Context, relationID sql.NullInt64) ([]*GetGroupNotifyByIDRow, error)
	GetGroupPinSettingsOrderByPinTime(ctx context.Context, accountID int64) ([]*GetGroupPinSettingsOrderByPinTimeRow, error)
	GetGroupRelationByID(ctx context.Context, id int64) (*GetGroupRelationByIDRow, error)
	GetGroupSettingsByName(ctx context.Context, arg *GetGroupSettingsByNameParams) ([]*GetGroupSettingsByNameRow, error)
	GetGroupShowSettingsOrderByShowTime(ctx context.Context, accountID int64) ([]*GetGroupShowSettingsOrderByShowTimeRow, error)
	GetMessageByID(ctx context.Context, id int64) (*GetMessageByIDRow, error)
	GetMsgsByContent(ctx context.Context, arg *GetMsgsByContentParams) ([]*GetMsgsByContentRow, error)
	GetMsgsByContentAndRelation(ctx context.Context, arg *GetMsgsByContentAndRelationParams) ([]*GetMsgsByContentAndRelationRow, error)
	GetMsgsByRelationIDAndTime(ctx context.Context, arg *GetMsgsByRelationIDAndTimeParams) ([]*GetMsgsByRelationIDAndTimeRow, error)
	GetPinMsgsByRelationID(ctx context.Context, arg *GetPinMsgsByRelationIDParams) ([]*GetPinMsgsByRelationIDRow, error)
	GetRelationIDByAccountID(ctx context.Context, accountID int64) (int64, error)
	GetRlyMsgsInfoByMsgID(ctx context.Context, arg *GetRlyMsgsInfoByMsgIDParams) ([]*GetRlyMsgsInfoByMsgIDRow, error)
	GetSettingByID(ctx context.Context, arg *GetSettingByIDParams) (*Setting, error)
	GetTopMsgByRelationID(ctx context.Context, relationID int64) (*GetTopMsgByRelationIDRow, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	OfferMsgsByAccountIDAndTime(ctx context.Context, arg *OfferMsgsByAccountIDAndTimeParams) ([]*OfferMsgsByAccountIDAndTimeRow, error)
	TransferIsLeaderFalse(ctx context.Context, arg *TransferIsLeaderFalseParams) error
	TransferIsLeaderTrue(ctx context.Context, arg *TransferIsLeaderTrueParams) error
	UpdateAccount(ctx context.Context, arg *UpdateAccountParams) error
	UpdateAccountAvatar(ctx context.Context, arg *UpdateAccountAvatarParams) error
	UpdateApplication(ctx context.Context, arg *UpdateApplicationParams) error
	UpdateGroupAvatar(ctx context.Context, arg *UpdateGroupAvatarParams) error
	UpdateGroupNotify(ctx context.Context, arg *UpdateGroupNotifyParams) (*UpdateGroupNotifyRow, error)
	UpdateGroupRelation(ctx context.Context, arg *UpdateGroupRelationParams) error
	UpdateMsgPin(ctx context.Context, arg *UpdateMsgPinParams) error
	UpdateMsgReads(ctx context.Context, arg *UpdateMsgReadsParams) ([]*UpdateMsgReadsRow, error)
	UpdateMsgRevoke(ctx context.Context, arg *UpdateMsgRevokeParams) error
	UpdateMsgTop(ctx context.Context, arg *UpdateMsgTopParams) error
	UpdateSettingDisturb(ctx context.Context, arg *UpdateSettingDisturbParams) error
	UpdateSettingLeader(ctx context.Context, arg *UpdateSettingLeaderParams) error
	UpdateSettingNickName(ctx context.Context, arg *UpdateSettingNickNameParams) error
	UpdateSettingPin(ctx context.Context, arg *UpdateSettingPinParams) error
	UpdateSettingShow(ctx context.Context, arg *UpdateSettingShowParams) error
	UpdateUser(ctx context.Context, arg *UpdateUserParams) error
}

var _ Querier = (*Queries)(nil)
