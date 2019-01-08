package User

import (
	"math/rand"
	"time"
	"fmt"

	"User/UserModel"
	"Util"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"

	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
)

func RegisterHanlder(ctx context.Context) {

	c := &UserModel.User{}
	name := Util.UnicodeEmojiCode(c.Name)
	account := Util.UnicodeEmojiCode(c.Account)
	databaseName := "memoryList"//默认
	databaseName += account//默认
	password := Util.UnicodeEmojiCode(c.Password)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	inviteCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
			//ceshi...
			name = "111"
			account = "1231233312"
			password = "12312333"
			databaseName = account
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
		_, _ = usersDB.Exec("CREATE TABLE IF NOT EXISTS users(name TEXT , databaseName TEXT , account TEXT , password TEXT,inviteCode TEXT)")
		// 查询单条数据
		row := usersDB.QueryRow("SELECT account FROM users WHERE account = ?", account)
		var searchAccount string
		err = row.Scan(&searchAccount) //遍历结果
		if err == nil && len(searchAccount) > 0 {//查到了
			fmt.Printf("手机号已经被注册\n")
			ctx.JSON(iris.Map{"errorCode": "501", "message": "手机号已经被注册"})
		} else {
			_, err := usersDB.Exec("INSERT INTO users(name, databaseName, account, password,inviteCode) VALUES(?,?,?,?,?)", name, databaseName, account, password,inviteCode) //插入数据
			if err == nil {
				fmt.Printf("注册成功,name:%s\ndatabaseName:%s\naccount:%s\npassword:%s\ninviteCode:%s\n", name, databaseName, account, password,inviteCode)
				ctx.JSON(iris.Map{"errorCode": "0", "message": "注册成功"})
			} else {
				fmt.Printf("注册失败:数据库插入失败\n")
				ctx.JSON(iris.Map{"errorCode": "500", "message": "注册失败"})
				Util.CheckError(err)
			}
		}
	}
}
