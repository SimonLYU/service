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

func RegisterHanlder(ctx context.Context) {

	c := &UserModel.User{}
	name := Util.UnicodeEmojiCode(c.Name)
	account := Util.UnicodeEmojiCode(c.Account)
	databaseName := account//默认
	password := Util.UnicodeEmojiCode(c.Password)

			//ceshi...
			// name = "111"
			// account = "12312333"
			// password = "12312333"
			// databaseName = account
			//ceshi over

	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		if len(password) == 0 || len(account) == 0 || len(name) == 0  {
			fmt.Printf("注册失败:不能为空\n")
			ctx.JSON(iris.Map{"errorCode": "500", "message": "注册失败:信息不完整"})
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
		_, _ = usersDB.Exec("CREATE TABLE IF NOT EXISTS users(name TEXT , databaseName TEXT , account TEXT , password TEXT)")
		// 查询单条数据
		row := usersDB.QueryRow("select * from users where account = ?", account)
		var searchAccount , searchPassword , searchName , searchDatabaseName string
		err = row.Scan(&searchName, &searchDatabaseName , &searchAccount,&searchPassword) //遍历结果
		if err == nil {
			fmt.Printf("手机号已经被注册\n")
			ctx.JSON(iris.Map{"errorCode": "501", "message": "手机号已经被注册"})
		} else {
			_, err := usersDB.Exec("insert into users(name, databaseName, account, password) values(?,?,?,?)", name, databaseName, account, password) //插入数据
			if err == nil {
				fmt.Printf("注册成功,name:%s\ndatabaseName:%s\naccount:%s\npassword:%s\n", name, databaseName, account, password)
				ctx.JSON(iris.Map{"errorCode": "0", "message": "注册成功"})
			} else {
				fmt.Printf("注册失败:数据库插入失败\n")
				ctx.JSON(iris.Map{"errorCode": "500", "message": "注册失败"})
			}
		}
	}
}
