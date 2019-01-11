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

func ChangeLinkDatabaseHandler(ctx context.Context) {

	c := &UserModel.User{}

	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		name := Util.UnicodeEmojiCode(c.Name)
		account := Util.UnicodeEmojiCode(c.Account)
		linkAccount := Util.UnicodeEmojiCode(c.LinkAccount)
		linkInviteCode := Util.UnicodeEmojiCode(c.LinkInviteCode)

		fullTimeString := time.Now().String()
		timeString := fullTimeString[:19]
		fmt.Printf("TIME:%s --> ", timeString)
		
		if len(linkAccount) == 0 || len(linkInviteCode) == 0 || len(account) == 0 {
			fmt.Printf("更改数据库失败:不能为空->linkAccount:%s,linkInviteCode:%s,account:%s\n",linkAccount,linkInviteCode,account)
			ctx.JSON(iris.Map{"errorCode": "500", "message": "变更失败:信息不完整", "databaseName": "", "linkName": ""})
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
		//初始化数据库
		_, _ = usersDB.Exec("CREATE TABLE IF NOT EXISTS users(name TEXT , databaseName TEXT , account TEXT , password TEXT,inviteCode TEXT,linkName TEXT)")
		row := usersDB.QueryRow("SELECT linkName,databaseName,inviteCode FROM users WHERE account = ?", linkAccount)
		var searchLinkName, searchLinkDatabase, searchLinkInviteCode string
		err = row.Scan(&searchLinkName, &searchLinkDatabase, &searchLinkInviteCode) //遍历结果
		if err == nil {                                                             //查到了
			if searchLinkInviteCode == linkInviteCode {
				var linkName, databaseName string

				if linkAccount == account {
					linkName = name
					databaseName = "dbn_"
					databaseName += account
				} else {
					linkName = searchLinkName
					databaseName = searchLinkDatabase
				}

				_, err = usersDB.Exec("update users set databaseName=?,linkName=? where account = ?", databaseName, linkName, account)
				if err == nil {
					fmt.Printf("link更新成功\n")
					ctx.JSON(iris.Map{"errorCode": "0", "message": "变更成功", "databaseName": databaseName, "linkName": linkName})
					return
				} else {
					fmt.Printf("link更新数据库失败->databaseName:%s,linkName:%s,linkAccount:%s\n", databaseName, linkName, linkAccount)
					ctx.JSON(iris.Map{"errorCode": "501", "message": "对方账号或邀请码填写不正确", "databaseName": "", "linkName": ""})
					return
				}
			} else {
				fmt.Printf("link邀请码不正确->linkAccount:%s,searchLinkInviteCode:%s,linkInviteCode:%s\n", linkAccount, searchLinkInviteCode, linkInviteCode)
				ctx.JSON(iris.Map{"errorCode": "501", "message": "对方账号或邀请码填写不正确", "databaseName": "", "linkName": ""})
				return
			}

		} else {
			fmt.Printf("link信息未查到->linkAccount:%s\n", linkAccount)
			ctx.JSON(iris.Map{"errorCode": "501", "message": "对方账号或邀请码填写不正确", "databaseName": "", "linkName": ""})
			return
		}
	}
}

func ChangeNameHandler(ctx context.Context) {

	c := &UserModel.User{}

	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		name := Util.UnicodeEmojiCode(c.Name)
		account := Util.UnicodeEmojiCode(c.Account)

		if len(name) == 0 || len(account) == 0 {
			fmt.Printf("更改数据库失败:不能为空\n")
			ctx.JSON(iris.Map{"errorCode": "500", "message": "变更失败:信息不完整", "name": ""})
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
		//初始化数据库
		_, _ = usersDB.Exec("CREATE TABLE IF NOT EXISTS users(name TEXT , databaseName TEXT , account TEXT , password TEXT,inviteCode TEXT,linkName TEXT)")
		_, err = usersDB.Exec("update users set name=? where account = ?", name, account)
		if err == nil {
			fmt.Printf("修改昵称成功->account:%s,name:%s\n", account, name)
			ctx.JSON(iris.Map{"errorCode": "0", "message": "修改成功", "name": name})
			return
		} else {
			fmt.Printf("修改昵称失败->account:%s,name:%s\n", account, name)
			ctx.JSON(iris.Map{"errorCode": "501", "message": "修改失败", "name": ""})
			return
		}
	}
}

func ChangePasswordHandler(ctx context.Context) {
	c := &UserModel.Password{}
	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		account := Util.UnicodeEmojiCode(c.Account)
		oldPassword := Util.UnicodeEmojiCode(c.OldPassword)
		newPassword := Util.UnicodeEmojiCode(c.NewPassword)

		if len(account) == 0 || len(oldPassword) == 0 || len(newPassword) == 0 {
			fmt.Printf("更改数据库失败:不能为空\n")
			ctx.JSON(iris.Map{"errorCode": "500", "message": "变更失败:信息不完整"})
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
		//初始化数据库
		_, _ = usersDB.Exec("CREATE TABLE IF NOT EXISTS users(name TEXT , databaseName TEXT , account TEXT , password TEXT,inviteCode TEXT,linkName TEXT)")
		row := usersDB.QueryRow("SELECT password FROM users WHERE account = ? AND password = ?", account, oldPassword)
		var searchPassword string
		err = row.Scan(&searchPassword) //遍历结果
		Util.CheckError(err)
		if err == nil { //查到了
			_, err = usersDB.Exec("UPDATE users SET password=? WHERE account = ?", newPassword, account)
			if err == nil {
				fmt.Printf("修改密码成功\n")
				ctx.JSON(iris.Map{"errorCode": "0", "message": "修改密码成功,请重新登录"})
				return
			} else {
				fmt.Printf("修改密码更新数据库失败->account:%s,oldPassword%s,newPassword:%s\n", account, oldPassword, newPassword)
				ctx.JSON(iris.Map{"errorCode": "501", "message": "修改密码失败"})
				return
			}
		} else {
			fmt.Printf("账号或密码不正确->account:%s,password:%s\n", account, oldPassword)
			ctx.JSON(iris.Map{"errorCode": "501", "message": "账号或密码不正确"})
			return
		}
	}
}
