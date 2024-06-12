package router

import (
	"ginchat/docs"
	"ginchat/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	//swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")

	//首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.GET("/chat", service.Chat)
	r.POST("/searchFriends", service.SearchFriends)

	//用户模块
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/getUserList", service.GetUserList)
	r.DELETE("/user/deleteUser/:id", service.DeleteUser)
	r.PUT("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)

	//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	//发送文件
	r.POST("/attach/upload", service.Upload)

	//关系操作
	r.POST("/contact/addFriend", service.AddFriend)
	r.POST("/contact/createCommunity", service.CreateCommunity)
	return r
}
