package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string // token
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     string
	HeartbeatTime string
	LoginOutTime  string `gorm:"columns:login_out_time" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

func FindUserByNameAndPwd(name, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and password = ?", name, password).First(&user)

	return user
}

func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", email).First(&user)
	return user
}

func CreateUser(user UserBasic) *gorm.DB {

	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {

	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {

	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, Password: user.Password, Phone: user.Phone, Email: user.Email})
}

func FindByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id = ?", id).First(&user)
	return user
}
