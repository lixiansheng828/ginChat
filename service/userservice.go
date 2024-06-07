package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetUserList
// @Summary 获取所有用户列表
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0,
		"message": data,
	})
}

// FindUserByNameAndPwd
// @Summary 根据用户名和密码查找用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	name := c.Request.FormValue("name")
	ver_password := c.Request.FormValue("password")
	user := models.FindUserByName(name)
	fmt.Println(ver_password, user.Password)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "该用户不存在",
		})
		return
	}
	ver_pwd := utils.MakePassword(ver_password, user.Salt)
	flag := utils.ValidPassword(ver_pwd, user.Password)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "密码有误",
		})
		return
	}
	data := models.FindUserByNameAndPwd(name, ver_pwd)
	//create token
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp_identity := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp_identity)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "登入成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param repassword formData string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {

	name := c.PostForm("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")
	if password != repassword {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "两次密码不一致!",
		})
		return
	}
	salt := fmt.Sprintf("%06d", rand.Int31())
	data := models.FindUserByName(name)
	if name == "" || password == "" || repassword == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名或密码为空!",
		})
		return
	}
	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名已注册",
		})
		return
	}
	user := models.UserBasic{}
	user.Password = utils.MakePassword(password, salt)
	user.Name = name
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "新增用户成功",
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @Param id path string true "用户ID"
// @Success 200 {string} json {"message":"删除成功"}
// @Router /user/deleteUser/{id} [delete]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	ID, _ := strconv.Atoi(c.Param("id")) // 通过路径参数获取ID
	user.ID = uint(ID)

	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "删除用户成功",
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json {"message":"修改成功"}
// @Router /user/updateUser [put]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}

	ID, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(ID)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	// 更新用户信息
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		c.JSON(200, gin.H{
			"code":  -1,
			"Error": err,
		})
		return
	}
	models.UpdateUser(user)

	c.JSON(200, gin.H{
		"code":    0,
		"message": "修改用户成功",
	})
}

// SearchFriends
// @Summary 查询好友列表
// @Tags 用户模块
// @param userId formData string false "id"
// @Success 200 {string} json {"message":"查询好友列表成功"}
// @Router /searchFriends [post]
func SearchFriends(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	fmt.Println("sa---", userId)
	users := models.SearchFriends(uint(userId))

	// c.JSON(200, gin.H{
	// 	"code":    0,
	// 	"message": "查询好友列表成功!",
	// 	"users":   users,
	// })

	utils.ResOKList(c.Writer, users, len(users))
}

// 防止跨域站点的伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	fmt.Println("SendMsg")
	MsgHandler(ws, c)

}
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println("failed to subscribe....", err)
		}
		tm := time.Now().Format("2006-1-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}

}
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)

}
