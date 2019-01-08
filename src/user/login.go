package User

import (
	"User/UserModel"
	"Util"
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"

	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
)

func LoginHanlder(ctx context.Context) {
	c := &UserModel.User{}
	account := c.Account
	password := c.Password


		//ceshi...
		// account = "123123"
		// password = "123123"
		//ceshi over

	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		if len(password) == 0 || len(account) == 0 {
			fmt.Printf("登录失败:不能为空\n")
			ctx.JSON(iris.Map{"errorCode": "502", "message": "登录失败:信息不完整"})
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
		_, _ = usersDB.Exec("CREATE TABLE IF NOT EXISTS users(name TEXT , databaseName TEXT , account TEXT , password TEXT)")
		// 查询多条数据
		row := usersDB.QueryRow("select * from users where account = ? and password = ?", account, password)
		var searchAccount , searchPassword , searchName , searchDatabaseName string
		err = row.Scan(&searchName, &searchDatabaseName , &searchAccount,&searchPassword) //遍历结果
		searchAccount = Util.UnicodeEmojiDecode(searchAccount)
		searchPassword = Util.UnicodeEmojiDecode(searchPassword)
		searchName = Util.UnicodeEmojiDecode(searchName)
		searchDatabaseName = Util.UnicodeEmojiDecode(searchDatabaseName)
		Util.CheckError(err)
		if err == nil{
			fmt.Printf("登录成功,name:%s\nacount:%s\ndatabaseName:%s\n",searchName,searchAccount,searchDatabaseName)
			ctx.JSON(iris.Map{"errorCode": "0", "message": "登录成功" , "name":searchName,"account":searchAccount,"databaseName":searchDatabaseName})
		}else{
			fmt.Printf("登录失败,用户名或密码错误\n")
			ctx.JSON(iris.Map{"errorCode": "502", "message": "登录失败:用户名或密码错误" , "name":"","account":"","databaseName":""})
		}	
	}
}
