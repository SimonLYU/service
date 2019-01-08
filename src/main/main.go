package main

import (
	"fmt"
	"time"
	"Utils/util"
	"user"

    "github.com/kataras/iris"
	"github.com/kataras/iris/context"

	"database/sql"
	_"github.com/Go-SQL-Driver/MySQL"
)

type Memory struct {
    Name  string
	MemoList []string
}

type Account struct {
	Name string
	AccountList []string
}

var staticMemoList []string 
var lastSetName string

var staticAccountList []string 
var lastUpdateAccountName string

func main() {

	memoryDB, openError := sql.Open("mysql", "root:simon.1314@/memoryList?charset=utf8")
	util.CheckError(openError)
	defer memoryDB.Close()//函数末关闭数据库

    // 设置最大打开的连接数，默认值为0表示不限制。
	memoryDB.SetMaxOpenConns(100)
	// 用于设置闲置的连接数。
	memoryDB.SetMaxIdleConns(50)
	memoryDB.Ping()
	//初始化DB
	_ , _ = memoryDB.Exec("CREATE TABLE IF NOT EXISTS memoList(memo TEXT)")
	_ , _ = memoryDB.Exec("CREATE TABLE IF NOT EXISTS accountList(info TEXT)")
	_ , _ = memoryDB.Exec("CREATE TABLE IF NOT EXISTS lastSetName(name TEXT)")
	_ , _ = memoryDB.Exec("CREATE TABLE IF NOT EXISTS lastUpdateAccountName(name TEXT)")

	// 查询多条数据
	rows, queryError := memoryDB.Query("SELECT memo FROM memoList")
	util.CheckError(queryError)
	// 对多条数据进行遍历
	var memo string
	for rows.Next() {
		scanError := rows.Scan(&memo)
		util.CheckError(scanError)
		//emoji表情解码
		fullTimeString := time.Now().String()
		timeString := fullTimeString[:19]
		fmt.Printf("TIME:%s --> 解码前数据:%s\n",timeString,memo)
		memo = util.UnicodeEmojiDecode(memo)
		fullTimeString = time.Now().String()
		timeString = fullTimeString[:19]
		fmt.Printf("TIME:%s --> 解码后数据:%s\n",timeString,memo)

		staticMemoList = append(staticMemoList, memo)
	}

	// 查询多条数据
	infoRows, queryInfoError := memoryDB.Query("SELECT info FROM accountList")
	util.CheckError(queryInfoError)
	// 对多条数据进行遍历
	var info string
	for infoRows.Next() {
		scanError := infoRows.Scan(&info)
		util.CheckError(scanError)
		//emoji表情解码
		fullTimeString := time.Now().String()
		timeString := fullTimeString[:19]
		fmt.Printf("TIME:%s --> 解码前数据:%s\n",timeString,info)
		info = util.UnicodeEmojiDecode(info)
		fullTimeString = time.Now().String()
		timeString = fullTimeString[:19]
		fmt.Printf("TIME:%s --> 解码后数据:%s\n",timeString,info)

		staticAccountList = append(staticAccountList, info)
	}

	nameRows, queryError := memoryDB.Query("SELECT name FROM lastSetName")
	util.CheckError(queryError)
	var name string
	for nameRows.Next() {
		scanError := nameRows.Scan(&name)
		util.CheckError(scanError)
		lastSetName = util.UnicodeEmojiDecode(name);
	}

	nameAccountRows, queryAccountError := memoryDB.Query("SELECT name FROM lastUpdateAccountName")
	util.CheckError(queryAccountError)
	var accountName string
	for nameAccountRows.Next() {
		scanError := nameAccountRows.Scan(&accountName)
		util.CheckError(scanError)
		lastUpdateAccountName = util.UnicodeEmojiDecode(accountName);
	}

	//开启服务器
	fullTimeString := time.Now().String()
	timeString := fullTimeString[:19]
	fmt.Printf("TIME:%s --> 服务器运行中...\n",timeString)
	postForMemoListApp()
}

func getMemoListHadnler(ctx context.Context) {
	nullMemoList := []string{}

	var currentMemoList []string
	if len(staticMemoList) <= 0 {
		currentMemoList = nullMemoList
	}else{
		currentMemoList = staticMemoList	
	}

	fullTimeString := time.Now().String()
	timeString := fullTimeString[:19]
	fmt.Printf("TIME:%s --> 购物车最后上报人:%s\t | 响应的列表:%s\n", timeString,lastSetName,currentMemoList)
	ctx.JSON(iris.Map{"name": lastSetName , "memoList" : currentMemoList})
}

func setMemoListHadnler(ctx context.Context) {
    c := &Memory{}
    if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
    } else {
		fullTimeString := time.Now().String()
		timeString := fullTimeString[:19]
        fmt.Printf("TIME:%s --> 购物车上报人:%s\t | 上报的列表:%#v\n",timeString,c.Name, c.MemoList)
		staticMemoList = c.MemoList
		lastSetName = c.Name
		//处理db
		db, openError := sql.Open("mysql", "root:simon.1314@/memoryList?charset=utf8")//本地数据库用户名root,密码:simon.1314
		util.CheckError(openError)
		defer db.Close()//函数末关闭数据库

		// 设置最大打开的连接数，默认值为0表示不限制。
		db.SetMaxOpenConns(100)
		// 用于设置闲置的连接数。
		db.SetMaxIdleConns(50)
		db.Ping()
		_, execError1 := db.Exec("DROP TABLE IF EXISTS memoList")
		util.CheckError(execError1)
		_, execError2 := db.Exec("CREATE TABLE memoList(memo TEXT)")
		util.CheckError(execError2)
	
		//这边事物内批量数据插入
		tx, dbError := db.Begin()
		util.CheckError(dbError)
		stmt, prepareError := tx.Prepare("INSERT memoList SET memo=?")
		util.CheckError(prepareError)
		for _, value := range  staticMemoList {
			//emoji表情转码
			value = util.UnicodeEmojiCode(value)

			_, err := stmt.Exec(value)
			if err != nil {
				fullTimeString := time.Now().String()
				timeString := fullTimeString[:19]
				fmt.Printf("TIME:%s --> 出现错误回滚，错误信息：%v\n",timeString, err)
				tx.Rollback()
			}
		}
		tx.Commit()

		//更新上报人
		_, execError3 := db.Exec("DROP TABLE IF EXISTS lastSetName")
		util.CheckError(execError3)
		_, execError4 := db.Exec("CREATE TABLE lastSetName(name TEXT)")
		util.CheckError(execError4)
		updateStmt,updateError := db.Prepare("INSERT INTO lastSetName(name) VALUES(?)")
		util.CheckError(updateError)
		_,execUpdateErr := updateStmt.Exec(util.UnicodeEmojiCode(lastSetName))
		util.CheckError(execUpdateErr)

		ctx.JSON(iris.Map{"name": lastSetName , "memoList" : staticMemoList})
    }
}

func updateAccountListHandler(ctx context.Context){

	c := &Account{}
    if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
    } else {
		fullTimeString := time.Now().String()
		timeString := fullTimeString[:19]
        fmt.Printf("TIME:%s --> 账本上报人:%s\t | 上报的列表:%#v\n",timeString,c.Name, c.AccountList)
		staticAccountList = c.AccountList
		lastUpdateAccountName = c.Name
		fmt.Printf("static = %s\n\n",c)
		//处理db
		db, openError := sql.Open("mysql", "root:simon.1314@/memoryList?charset=utf8")//本地数据库用户名root,密码:simon.1314
		util.CheckError(openError)
		defer db.Close()//函数末关闭数据库

		// 设置最大打开的连接数，默认值为0表示不限制。
		db.SetMaxOpenConns(100)
		// 用于设置闲置的连接数。
		db.SetMaxIdleConns(50)
		db.Ping()

		_, execError1 := db.Exec("DROP TABLE IF EXISTS accountList")
		util.CheckError(execError1)
		_, execError2 := db.Exec("CREATE TABLE accountList(info TEXT)")
		util.CheckError(execError2)
	
		//这边事物内批量数据插入
		tx, dbError := db.Begin()
		util.CheckError(dbError)
		stmt, prepareError := tx.Prepare("INSERT INTO accountList(info) values(?)")
		util.CheckError(prepareError)
		for _, value := range  staticAccountList {
			value = util.UnicodeEmojiCode(value)
			_, err := stmt.Exec(value)
			if err != nil {
				fullTimeString := time.Now().String()
				timeString := fullTimeString[:19]
				fmt.Printf("TIME:%s --> 出现错误回滚，错误信息：%v\n",timeString, err)
				tx.Rollback()
			}
		}
		tx.Commit()

		//更新上报人
		_, execError3 := db.Exec("DROP TABLE IF EXISTS lastUpdateAccountName")
		util.CheckError(execError3)
		_, execError4 := db.Exec("CREATE TABLE lastUpdateAccountName(name TEXT)")
		util.CheckError(execError4)
		updateStmt,updateError := db.Prepare("INSERT INTO lastUpdateAccountName(name) VALUES(?)")
		util.CheckError(updateError)
		_,execUpdateErr := updateStmt.Exec(util.UnicodeEmojiCode(lastUpdateAccountName))
		util.CheckError(execUpdateErr)

		ctx.JSON(iris.Map{"name": lastUpdateAccountName , "accountList" : staticAccountList})
    }
}

func getAccountListHadnler(ctx context.Context) {
	nullAccountList := []string{}

	var currentAccountList []string
	if len(staticAccountList) <= 0 {
		currentAccountList = nullAccountList
	}else{
		currentAccountList = staticAccountList	
	}

	fullTimeString := time.Now().String()
	timeString := fullTimeString[:19]
	fmt.Printf("TIME:%s --> 账本最后上报人:%s\t | 响应的列表:%s\n", timeString,lastUpdateAccountName,currentAccountList)
	ctx.JSON(iris.Map{"name": lastUpdateAccountName , "accountList" : currentAccountList})
}

func postForMemoListApp(){
	app := iris.New()
	//获取备忘录列表
	app.Post("/getMemoList", getMemoListHadnler)
	//全量更新备忘录列表
	app.Post("/setMemoList", setMemoListHadnler)
	//全量更新账单列表
	app.Post("/updateAccountList", updateAccountListHandler)
	//获取账单列表
	app.Post("/getAccountList", getAccountListHadnler)
	//登录
	app.Post("/login", user.LoginHanlder)
	//注册
	app.Post("/register", user.RegisterHanlder)
    app.Run(iris.Addr(":8080"))
}

