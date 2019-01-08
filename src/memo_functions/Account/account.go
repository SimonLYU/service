package Account

import (
	"fmt"
	"time"

	"Util"
	"memo_functions/Account/AccountModel"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"

	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
)

func UpdateAccountListHandler(ctx context.Context) {

	c := &AccountModel.Account{}
	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		fullTimeString := time.Now().String()
		timeString := fullTimeString[:19]
		fmt.Printf("TIME:%s --> 账本上报人:%s\t | 上报数据库:%s | 上报的列表:%#v\n", timeString, c.Name, c.DatabaseName, c.AccountList)
		staticAccountList := c.AccountList
		lastUpdateAccountName := c.Name
		databaseName := c.DatabaseName
		databaseName = "accountList" //ceshi...
		if len(databaseName) <= 0 {
			fmt.Printf("databaseName为空!")
			return
		}
		//处理db
		db, openError := sql.Open("mysql", "root:simon.1314@/account?charset=utf8") //本地数据库用户名root,密码:simon.1314
		Util.CheckError(openError)
		defer db.Close() //函数末关闭数据库

		// 设置最大打开的连接数，默认值为0表示不限制。
		db.SetMaxOpenConns(100)
		// 用于设置闲置的连接数。
		db.SetMaxIdleConns(50)
		db.Ping()

		drop := "DROP TABLE IF EXISTS "
		drop += databaseName
		_, execError1 := db.Exec(drop)
		Util.CheckError(execError1)
		create := "CREATE TABLE "
		create += databaseName
		create += "(info TEXT , name TEXT)"
		_, execError2 := db.Exec(create)
		Util.CheckError(execError2)

		//这边事物内批量数据插入
		tx, dbError := db.Begin()
		Util.CheckError(dbError)
		insert := "INSERT INTO "
		insert += databaseName
		insert += "(info,name) values(?,?)"
		stmt, prepareError := tx.Prepare(insert)
		Util.CheckError(prepareError)
		for _, value := range staticAccountList {
			value = Util.UnicodeEmojiCode(value)
			_, err := stmt.Exec(value, lastUpdateAccountName)
			if err != nil {
				fullTimeString := time.Now().String()
				timeString := fullTimeString[:19]
				fmt.Printf("TIME:%s --> 出现错误回滚，错误信息：%v\n", timeString, err)
				tx.Rollback()
			}
		}
		tx.Commit()

		ctx.JSON(iris.Map{"name": lastUpdateAccountName, "accountList": staticAccountList})
	}
}

func GetAccountListHadnler(ctx context.Context) {

	c := &AccountModel.GetAccount{}
	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {

		databaseName := c.DatabaseName
		databaseName = "accountList" //ceshi...
		if len(databaseName) <= 0 {
			fmt.Printf("databaseName为空!")
			return
		}

		memoryDB, openError := sql.Open("mysql", "root:simon.1314@/account?charset=utf8")
		Util.CheckError(openError)
		defer memoryDB.Close() //函数末关闭数据库
		// 设置最大打开的连接数，默认值为0表示不限制。
		memoryDB.SetMaxOpenConns(100)
		// 用于设置闲置的连接数。
		memoryDB.SetMaxIdleConns(50)
		memoryDB.Ping()

		//初始化数据库
		create := "CREATE TABLE "
		create += databaseName
		create += "(info TEXT , name TEXT)"
		_, _ = memoryDB.Exec(create)

		var staticAccountList []string
		var lastUpdateAccountName string
		// 查询多条数据
		selectExec := "SELECT info,name FROM "
		selectExec += databaseName
		infoRows, queryInfoError := memoryDB.Query(selectExec)
		Util.CheckError(queryInfoError)
		// 对多条数据进行遍历
		var info, name string
		for infoRows.Next() {
			scanError := infoRows.Scan(&info, &name)
			Util.CheckError(scanError)
			//emoji表情解码
			fullTimeString := time.Now().String()
			timeString := fullTimeString[:19]
			fmt.Printf("TIME:%s --> 解码前数据:%s\n", timeString, info)
			info = Util.UnicodeEmojiDecode(info)
			fullTimeString = time.Now().String()
			timeString = fullTimeString[:19]
			fmt.Printf("TIME:%s --> 解码后数据:%s\n", timeString, info)

			staticAccountList = append(staticAccountList, info)
		}
		lastUpdateAccountName = Util.UnicodeEmojiDecode(name)

		nullAccountList := []string{}
		var currentAccountList []string
		if len(staticAccountList) <= 0 {
			currentAccountList = nullAccountList
		} else {
			currentAccountList = staticAccountList
		}

		fullTimeString := time.Now().String()
		timeString := fullTimeString[:19]
		fmt.Printf("TIME:%s --> 账本最后上报人:%s\t | 响应的列表:%s\n", timeString, lastUpdateAccountName, currentAccountList)
		ctx.JSON(iris.Map{"name": lastUpdateAccountName, "accountList": currentAccountList})
	}
}
