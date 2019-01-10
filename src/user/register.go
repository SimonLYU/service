package User

import (
	"fmt"
	"math/rand"
	"time"

	"User/UserModel"
	"Util"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"

	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
)

func RegisterHanlder(ctx context.Context) {

	c := &UserModel.User{}

	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		name := Util.UnicodeEmojiCode(c.Name)
		account := Util.UnicodeEmojiCode(c.Account)
		password := Util.UnicodeEmojiCode(c.Password)
		linkAccount := Util.UnicodeEmojiCode(c.LinkAccount)
		linkInviteCode := Util.UnicodeEmojiCode(c.LinkInviteCode)
		fmt.Printf("link info--%s,%s\nregisterInfo--%s,%s,%s\n", linkAccount, linkInviteCode,name,account,password)
		//默认根据account生成本人的数据库
		databaseName := "dbn_"
		databaseName += account
		//默认使用自己的db,link名字为自己
		linkName := name
		//随机生成邀请码
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		inviteCode := Util.UnicodeEmojiCode(fmt.Sprintf("%06v", rnd.Int31n(1000000)))

		if len(password) == 0 || len(account) == 0 || len(name) == 0 {
			fmt.Printf("注册失败:不能为空\n")
			ctx.JSON(iris.Map{"errorCode": "500", "message": "注册失败:信息不完整", "name": "", "account": "", "databaseName": "", "inviteCode": "", "linkName": ""})
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
		// 查询单条数据
		row := usersDB.QueryRow("SELECT account FROM users WHERE account = ?", account)
		var searchAccount string
		err = row.Scan(&searchAccount)            //遍历结果
		if err == nil && len(searchAccount) > 0 { //查到了
			fmt.Printf("手机号已经被注册\n")
			ctx.JSON(iris.Map{"errorCode": "501", "message": "手机号已经被注册", "name": "", "account": "", "databaseName": "", "inviteCode": "", "linkName": ""})
			return
		} else {
			if len(linkAccount) != 0 && len(linkInviteCode) != 0 {
				row := usersDB.QueryRow("SELECT linkName,databaseName,inviteCode FROM users WHERE account = ?", linkAccount)
				var  searchLinkName, searchLinkDatabase, searchLinkInviteCode string
				err = row.Scan(&searchLinkName, &searchLinkDatabase, &searchLinkInviteCode) //遍历结果
				if err == nil {                                                                                 //查到了
					if searchLinkInviteCode == linkInviteCode {
						linkName = searchLinkName
						databaseName = searchLinkDatabase
					} else {
						fmt.Printf("link邀请码不正确,name:%s\ndatabaseName:%s\naccount:%s\npassword:%s\ninviteCode:%s\nlinkName:%s\n", name, databaseName, account, password, inviteCode, linkName)
						ctx.JSON(iris.Map{"errorCode": "501", "message": "对方账号或邀请码填写不正确", "name": "", "account": "", "databaseName": "", "inviteCode": "", "linkName": ""})
						return
					}

				} else {
					fmt.Printf("link信息未查到,name:%s\ndatabaseName:%s\naccount:%s\npassword:%s\ninviteCode:%s\nlinkName:%s\n", name, databaseName, account, password, inviteCode, linkName)
					ctx.JSON(iris.Map{"errorCode": "501", "message": "对方账号或邀请码填写不正确", "name": "", "account": "", "databaseName": "", "inviteCode": "", "linkName": ""})
					return
				}
			} else if (len(linkAccount) == 0 && len(linkInviteCode) != 0) || len(linkAccount) != 0 && len(linkInviteCode) == 0 {
				fmt.Printf("link填写不完整\n")
				ctx.JSON(iris.Map{"errorCode": "501", "message": "对方账号或邀请码填写不正确", "name": "", "account": "", "databaseName": "", "inviteCode": "", "linkName": ""})
				return
			}

			_, err := usersDB.Exec("INSERT INTO users(name, databaseName, account, password,inviteCode,linkName) VALUES(?,?,?,?,?,?)", name, databaseName, account, password, inviteCode, linkName) //插入数据
			if err == nil {
				fmt.Printf("注册成功,name:%s\ndatabaseName:%s\naccount:%s\npassword:%s\ninviteCode:%s\nlinkName:%s\n", name, databaseName, account, password, inviteCode, linkName)
				ctx.JSON(iris.Map{"errorCode": "0", "message": "注册成功,请登录", "name": name, "account": account, "databaseName": databaseName, "inviteCode": inviteCode, "linkName": linkName})
				return
			} else {
				fmt.Printf("注册失败,name:%s\ndatabaseName:%s\naccount:%s\npassword:%s\ninviteCode:%s\nlinkName:%s\n", name, databaseName, account, password, inviteCode, linkName)
				ctx.JSON(iris.Map{"errorCode": "500", "message": "注册失败", "name": "", "account": "", "databaseName": "", "inviteCode": "", "linkName": ""})
				return
			}
		}
	}
}
