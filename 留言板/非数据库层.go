package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func main(){
	router:=gin.Default()

	router.POST("/snow/registe", registe)//注册
	router.POST("/snow/information/registe",informationregiste)
	router.POST("/snow/information/selete",informationselete)
	router.POST("/snow/login",func(c *gin.Context){//登录
		username := c.Query("username")
		password := c.Query("password")
		if Loginll(username, password){ c.JSON(200, gin.H{"登录成功！": c.Query("username"),})
		}else {c.JSON(400, gin.H{"账号或密码错误": "!",})}
	})

	router.POST("/snow/findmessage",func(c *gin.Context){//留言
		pid:=c.Query("pid")
		pidNew, _ :=strconv.Atoi(pid)
		res:=FindMessageByPid(pidNew)
		c.JSON(200,JsonNested(res))
	})

	router.POST("/snow/insert",insert)
	router.DELETE("/snow/delete",Delete)

	//router.POST("/snow/test",func(c *gin.Context){
	//	c.JSON(200,gin.H{"id":SelectIdMwx()})
	//})

	_ = router.Run(":8080")
}

func Delete(c *gin.Context){
	id:=c.Query("id")
	iid,_:=strconv.Atoi(id)
    fmt.Println(iid)
	if DeleteMesage(iid){
		c.JSON(200, gin.H{"Data": "删除成功"})
	}else {
		c.JSON(500,gin.H{"status": http.StatusInternalServerError,"message":"数据库删除失败"})

	}
}

func insert(c *gin.Context){
	user_id := c.Query("user_id")
	user_idNew, _ :=strconv.Atoi(user_id)
	message := c.Query("message")
	pid := c.Query("pid")
	pidNew, _ :=strconv.Atoi(pid)
	//fmt.Println("user_idNew:",user_idNew,"message:",message,"pidNew:",pidNew)
    id := SelectIdMax()+1
	if MessageInsert(pidNew,id,user_idNew,message){
		if pidNew == 0 {
			c.JSON(200, gin.H{"status": http.StatusOK, "message": "在别人主页留言成功"})
		} else{
		c.JSON(200, gin.H{"status": http.StatusOK, "message": "回复别人留言成功"})
		}
	}else {
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"数据库Insert报错"})
	}
}

func registe(c *gin.Context){//注册
	username := c.Query("username")
	password := c.Query("password")
	id :=c.Query("id")
	if UserExist(username){
	if UserSignup(username,password,id){
		c.JSON(200, gin.H{"status": http.StatusOK, "message": "注册成功"})
	}else {
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"数据库Insert报错"})
	}}else{
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"该用户已注册"})
	}
}

func informationregiste(c *gin.Context){
	username := c.Query("username")
	gender := c.Query("gender")
	introduce := c.Query("introduce")
	address := c.Query("address")
	industry := c.Query("industry")
	occupation := c.Query("occupation")
	education := c.Query("education")
	if UserExist(username){
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"该用户未注册"})
	}else{
	if InformationSignup(username,gender,introduce,address,industry,occupation,education){
		c.JSON(200, gin.H{"status": http.StatusOK, "message": "注册成功"})
	}else {
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"message":"数据库Insert报错"})
	}}
}

func informationselete(c *gin.Context){
	username := c.Query("username")
	gender,introduce,address,industry,occupation,education:=InformationSelect(username)
		c.JSON(200, gin.H{"status": http.StatusOK, "gender": gender, "introduce": introduce, "address": address, "industry": industry, "occupation": occupation, "education": education})
}

type Message struct {
	UserId int
	Message string
	ChildMessage *[]Message
}

//var order = 0

func JsonNested(messageSlice []Message) []gin.H {
	//order++
	var messageJsons []gin.H
	//fmt.Printf("第%d层开始", order)
	//fmt.Println()
	var messageJson gin.H
	for _, messages := range messageSlice {
		//fmt.Println("分解过程", messages)
		message := *messages.ChildMessage
		//fmt.Println("分解过程的的子留言", message)
		if messages.ChildMessage != nil {
			messageJson = gin.H{
				"user_id":         messages.UserId,
				"message":         messages.Message,
				"ChildrenMessage": JsonNested(message),
			}
		} else {
			messageJson = gin.H{
				"user_id": messages.UserId,
				"message": messages.Message,
				"ChildrenMessage":"null",
			}
		}
		messageJsons = append(messageJsons, messageJson)
	}
	//fmt.Printf("第%d层结束。", order)
	//fmt.Println()
	//order--
	return messageJsons
}
