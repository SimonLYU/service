package Memory

import (
	"fmt"
	"time"

	"Util"
	"memo_functions/Memory/MemoryModel"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"

	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
)

func GetMemoListHadnler(ctx context.Context) {
	c := &MemoryModel.GetMemory{}
	fullTimeString := time.Now().String()
	timeString := fullTimeString[:19]
	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		databaseName := c.DatabaseName
		if len(databaseName) <= 0 {
			fmt.Printf("GetMemoList databaseName为空!\n")
			return
		}

		memoryDB, openError := sql.Open("mysql", "root:simon.1314@/memory?charset=utf8")
		Util.CheckError(openError)
		defer memoryDB.Close() //函数末关闭数据库
		// 设置最大打开的连接数，默认值为0表示不限制。
		memoryDB.SetMaxOpenConns(100)
		// 用于设置闲置的连接数。
		memoryDB.SetMaxIdleConns(50)
		memoryDB.Ping()
		//初始化数据库
		create := "CREATE TABLE IF NOT EXISTS "
		create += databaseName
		create += "(memo TEXT , name TEXT)"
		_, _ = memoryDB.Exec(create)

		var staticMemoList []string
		var lastSetName string
		// 查询多条数据
		selectExec := "SELECT memo,name FROM "
		selectExec += databaseName
		rows, queryError := memoryDB.Query(selectExec)
		Util.CheckError(queryError)
		// 对多条数据进行遍历
		var memo, name string
		for rows.Next() {
			scanError := rows.Scan(&memo, &name)
			Util.CheckError(scanError)
			//emoji表情解码
			memo = Util.UnicodeEmojiDecode(memo)
			// fmt.Printf("TIME:%s --> 解码后数据:%s\n", timeString, memo)

			staticMemoList = append(staticMemoList, memo)
		}

		lastSetName = Util.UnicodeEmojiDecode(name)

		nullMemoList := []string{}
		var currentMemoList []string
		if len(staticMemoList) <= 0 {
			currentMemoList = nullMemoList
		} else {
			currentMemoList = staticMemoList
		}

		fmt.Printf("TIME:%s -->GetMemoList 购物车最后上报人:%s\t | 响应的列表:%s\n", timeString, lastSetName, currentMemoList)
		ctx.JSON(iris.Map{"name": lastSetName, "memoList": currentMemoList})
	}
}

func UpdateMemoListHadnler(ctx context.Context) {
	c := &MemoryModel.Memory{}
	fullTimeString := time.Now().String()
	timeString := fullTimeString[:19]
	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("TIME:%s -->UpdateMemo 购物车上报人:%s\t | 上报的列表:%#v\n", timeString, c.Name, c.MemoList)
		staticMemoList := c.MemoList
		lastSetName := c.Name
		databaseName := c.DatabaseName
		if len(databaseName) <= 0 {
			fmt.Printf("UpdateMemo databaseName为空!\n")
			return
		}
		//处理db
		db, openError := sql.Open("mysql", "root:simon.1314@/memory?charset=utf8") //本地数据库用户名root,密码:simon.1314
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
		create += "(memo TEXT , name TEXT)"
		_, execError2 := db.Exec(create)
		Util.CheckError(execError2)

		//这边事物内批量数据插入
		tx, dbError := db.Begin()
		Util.CheckError(dbError)
		insert := "INSERT INTO "
		insert += databaseName
		insert += "(memo,name) values(?,?)"
		stmt, prepareError := tx.Prepare(insert)
		Util.CheckError(prepareError)
		for _, value := range staticMemoList {
			//emoji表情转码
			value = Util.UnicodeEmojiCode(value)

			_, err := stmt.Exec(value, lastSetName)
			if err != nil {
				fmt.Printf("TIME:%s -->UpdateMemo 出现错误回滚，错误信息：%v\n", timeString, err)
				tx.Rollback()
			}
		}
		tx.Commit()

		ctx.JSON(iris.Map{"name": lastSetName, "memoList": staticMemoList})
	}
}

func InsertMemoListHadnler(ctx context.Context) {
	c := &MemoryModel.SingleMemory{}
	fullTimeString := time.Now().String()
	timeString := fullTimeString[:19]
	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("TIME:%s -->InsertMemoList 购物车上报人:%s\t | 上报的单条:%#v\n", timeString, c.Name, c.Memo)
		staticMemo := c.Memo
		lastSetName := c.Name
		databaseName := c.DatabaseName
		if len(databaseName) <= 0 {
			fmt.Printf("InsertMemoList databaseName为空!\n")
			return
		}
		//处理db
		db, openError := sql.Open("mysql", "root:simon.1314@/memory?charset=utf8") //本地数据库用户名root,密码:simon.1314
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
		create += "(memo TEXT , name TEXT)"
		_, execError2 := db.Exec(create)
		Util.CheckError(execError2)

		//这边事物内批量数据插入
		insert := "INSERT INTO "
		insert += databaseName
		insert += "(memo,name) values(?,?)"

		_,err = db.Exec(insert,staticMemo,lastSetName)
		Util.CheckError(err)

		var staticMemoList []string
		// 查询多条数据
		selectExec := "SELECT memo,name FROM "
		selectExec += databaseName
		rows, queryError := db.Query(selectExec)
		Util.CheckError(queryError)
		// 对多条数据进行遍历
		var memo, name string
		for rows.Next() {
			scanError := rows.Scan(&memo, &name)
			Util.CheckError(scanError)
			//emoji表情解码
			memo = Util.UnicodeEmojiDecode(memo)
			// fmt.Printf("TIME:%s -->InsertMemoList 解码后数据:%s\n", timeString, memo)

			staticMemoList = append(staticMemoList, memo)
		}

		lastSetName = Util.UnicodeEmojiDecode(name)

		nullMemoList := []string{}
		var currentMemoList []string
		if len(staticMemoList) <= 0 {
			currentMemoList = nullMemoList
		} else {
			currentMemoList = staticMemoList
		}
		fmt.Printf("TIME:%s -->InsertMemoList 返回的数据:%s\n", timeString, currentMemoList)
		ctx.JSON(iris.Map{"name": lastSetName, "memoList": currentMemoList})
	}
}