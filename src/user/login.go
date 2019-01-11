package User

import (
	"User/UserModel"
	"Util"
	"fmt"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"

	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
)

func LoginHanlder(ctx context.Context) {
	c := &UserModel.User{}

	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		account := c.Account
		password := c.Password
		fullTimeString := time.Now().String()
		timeString := fullTimeString[:19]
		fmt.Printf("TIME:%s --> ", timeString)
		fmt.Printf("登录请求->\nacount:%s\npassword:%s\n", account, password)
		fmt.Printf("TIME:%s --> ", timeString)
		if len(account) == 0 || len(password) == 0 {
			fmt.Printf("信息不完整->account:%s,password:%s\n",account,password)
			ctx.JSON(iris.Map{"errorCode": "502", "message": "登录失败:信息不完整", "name": "", "account": "", "databaseName": "","inviteCode":"","linkName":""})
			return
		}
		usersDB, openError := sql.Open("mysql", "root:simon.1314@/users?charset=utf8")
		Util.CheckError(openError)
		defer usersDB.Close() //函数末关闭数据库
		// 设置最大打开的连接数，默认值为0表示不限制。
		usersDB.SetMaxOpenConns(100)
		// 用于设置闲置的连接数。
		usersDB.SetMaxIdleConns(50)
		usersDB.Ping()
		// 初始化数据库
		_, _ = usersDB.Exec("CREATE TABLE IF NOT EXISTS users(name TEXT , databaseName TEXT , account TEXT , password TEXT,inviteCode TEXT,linkName TEXT)")
		// 查询单条数据
		row := usersDB.QueryRow("SELECT name,databaseName,account,inviteCode,linkName FROM users WHERE account = ? AND password = ?", account, password)
		var searchAccount, searchName, searchDatabaseName,searchInviteCode,searchLinkName string
		err = row.Scan(&searchName, &searchDatabaseName, &searchAccount , &searchInviteCode,&searchLinkName) //遍历结果
		searchAccount = Util.UnicodeEmojiDecode(searchAccount)
		searchName = Util.UnicodeEmojiDecode(searchName)
		searchDatabaseName = Util.UnicodeEmojiDecode(searchDatabaseName)
		searchInviteCode = Util.UnicodeEmojiDecode(searchInviteCode)
		searchLinkName = Util.UnicodeEmojiDecode(searchLinkName)
		// Util.CheckError(err)
		if err == nil {
			fmt.Printf("登录成功,name:%s\n\tacount:%s\n\tdatabaseName:%s\n\tinviteCode:%s\n\tlinkName:%s\n", searchName, searchAccount, searchDatabaseName,searchInviteCode,searchLinkName)
			ctx.JSON(iris.Map{"errorCode": "0", "message": "登录成功", "name": searchName, "account": searchAccount, "databaseName": searchDatabaseName,"inviteCode":searchInviteCode,"linkName":searchLinkName})
		} else {
			fmt.Printf("登录失败,用户名或密码错误->account:%s,password%s\n",account,password)
			ctx.JSON(iris.Map{"errorCode": "502", "message": "登录失败:用户名或密码错误", "name": "", "account": "", "databaseName": "","inviteCode":"","linkName":""})
		}
	}
}