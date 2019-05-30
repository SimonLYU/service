package main

import (
	"fmt"
	"time"

	"User"
	"memo_functions/Account"
	"memo_functions/Memory"

    "github.com/kataras/iris"
)

var staticMemoList []string 
var lastSetName string

var staticAccountList []string 
var lastUpdateAccountName string

func main() {

	//开启服务器
	fullTimeString := time.Now().String()
	timeString := fullTimeString[:19]
	fmt.Printf("TIME:%s --> 服务器运行中...\n",timeString)
	postForMemoListApp()
}

func postForMemoListApp(){
	app := iris.New()
	//获取备忘录列表
	app.Post("/getMemoList", Memory.GetMemoListHadnler)
	//全量更新备忘录列表
	app.Post("/setMemoList", Memory.UpdateMemoListHadnler)
	//插入一条备忘(暂时不用)
	app.Post("/insertToMemoList", Memory.InsertMemoListHadnler)
	//全量更新账单列表
	app.Post("/updateAccountList", Account.UpdateAccountListHandler)
	//获取账单列表
	app.Post("/getAccountList", Account.GetAccountListHadnler)
	//插入一条账单列表
	app.Post("/insertAccountList", Account.InsertAccountListHandler)
	//删除一条账单列表
	app.Post("/deleteAccountList", Account.DeleteAccountListHandler)
	//登录
	app.Post("/login", User.LoginHanlder)
	//注册
	app.Post("/register", User.RegisterHanlder)
	//管更数据库
	app.Post("/changeLinkDatabase", User.ChangeLinkDatabaseHandler)
	//修改昵称
	app.Post("/changeUserInfo", User.ChangeNameHandler)
	//修改密码
	app.Post("/changePassword", User.ChangePasswordHandler)

    app.Run(iris.Addr(":8081"))
}