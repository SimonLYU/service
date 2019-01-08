package Account

import (
	"fmt"
	"time"

	"memo_functions/Account/AccountModel"
	"Util"

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
		fmt.Printf("TIME:%s --> 账本上报人:%s\t | 上报的列表:%#v\n", timeString, c.Name, c.AccountList)
		staticAccountList := c.AccountList
		lastUpdateAccountName := c.Name
		fmt.Printf("static = %s\n\n", c)
		//处理db
		db, openError := sql.Open("mysql", "root:simon.1314@/memoryList?charset=utf8") //本地数据库用户名root,密码:simon.1314
		Util.CheckError(openError)
		defer db.Close() //函数末关闭数据库

		// 设置最大打开的连接数，默认值为0表示不限制。
		db.SetMaxOpenConns(100)
		// 用于设置闲置的连接数。
		db.SetMaxIdleConns(50)
		db.Ping()

		_, execError1 := db.Exec("DROP TABLE IF EXISTS accountList")
		Util.CheckError(execError1)
		_, execError2 := db.Exec("CREATE TABLE accountList(info TEXT)")
		Util.CheckError(execError2)

		//这边事物内批量数据插入
		tx, dbError := db.Begin()
		Util.CheckError(dbError)
		stmt, prepareError := tx.Prepare("INSERT INTO accountList(info) values(?)")
		Util.CheckError(prepareError)
		for _, value := range staticAccountList {
			value = Util.UnicodeEmojiCode(value)
			_, err := stmt.Exec(value)
			if err != nil {
				fullTimeString := time.Now().String()
				timeString := fullTimeString[:19]
				fmt.Printf("TIME:%s --> 出现错误回滚，错误信息：%v\n", timeString, err)
				tx.Rollback()
			}
		}
		tx.Commit()

		//更新上报人
		_, execError3 := db.Exec("DROP TABLE IF EXISTS lastUpdateAccountName")
		Util.CheckError(execError3)
		_, execError4 := db.Exec("CREATE TABLE lastUpdateAccountName(name TEXT)")
		Util.CheckError(execError4)
		updateStmt, updateError := db.Prepare("INSERT INTO lastUpdateAccountName(name) VALUES(?)")
		Util.CheckError(updateError)
		_, execUpdateErr := updateStmt.Exec(Util.UnicodeEmojiCode(lastUpdateAccountName))
		Util.CheckError(execUpdateErr)

		ctx.JSON(iris.Map{"name": lastUpdateAccountName, "accountList": staticAccountList})
	}
}

func GetAccountListHadnler(ctx context.Context) {
	memoryDB, openError := sql.Open("mysql", "root:simon.1314@/memoryList?charset=utf8")
	Util.CheckError(openError)
	defer memoryDB.Close() //函数末关闭数据库
	// 设置最大打开的连接数，默认值为0表示不限制。
	memoryDB.SetMaxOpenConns(100)
	// 用于设置闲置的连接数。
	memoryDB.SetMaxIdleConns(50)
	memoryDB.Ping()

	//初始化数据库
	_ , _ = memoryDB.Exec("CREATE TABLE IF NOT EXISTS accountList(info TEXT)")
	_ , _ = memoryDB.Exec("CREATE TABLE IF NOT EXISTS lastUpdateAccountName(name TEXT)")

	var staticAccountList []string
	var lastUpdateAccountName string
	// 查询多条数据
	infoRows, queryInfoError := memoryDB.Query("SELECT info FROM accountList")
	Util.CheckError(queryInfoError)
	// 对多条数据进行遍历
	var info string
	for infoRows.Next() {
		scanError := infoRows.Scan(&info)
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

	nameAccountRows, queryAccountError := memoryDB.Query("SELECT name FROM lastUpdateAccountName")
	Util.CheckError(queryAccountError)
	var accountName string
	for nameAccountRows.Next() {
		scanError := nameAccountRows.Scan(&accountName)
		Util.CheckError(scanError)
		lastUpdateAccountName = Util.UnicodeEmojiDecode(accountName)
	}

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
