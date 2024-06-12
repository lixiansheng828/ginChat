package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint //主
	TargetId uint //客
	Type     int  //对应的类型 1好友 2群 3xx
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

// 查找好友列表
func SearchFriends(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id=?", userId).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(v)
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}

func AddFriend(userId uint, targetId uint) (int, string) {
	if targetId != 0 {
		if targetId == userId {
			return -1, "不允许添加自己"
		}
		contact := Contact{}
		utils.DB.Where("owner_id = ? and target_id = ?", userId, targetId).Find(&contact)
		if contact.ID != 0 {
			return -1, "重复添加"
		}
		user := FindByID(targetId)

		if user.Salt != "" {
			tx := utils.DB.Begin()
			// 事务一旦开始，不论什么异常最终都会Rollback
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()
			var contactMy Contact = Contact{
				OwnerId:  userId,
				TargetId: targetId,
				Type:     1,
			}
			var contactYou Contact = Contact{
				OwnerId:  targetId,
				TargetId: userId,
				Type:     1,
			}
			if err := utils.DB.Create(&contactMy).Error; err != nil {
				tx.Rollback()
				return -1, "添加失败"
			}
			if err := utils.DB.Create(&contactYou).Error; err != nil {
				tx.Rollback()
				return -1, "添加失败"
			}
			tx.Commit()
			return 0, "添加成功"
		}
	}

	return -1, "查无此人"
}
