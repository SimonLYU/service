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
	memoryDB, openError := sql.Open("mysql", "root:simon.1314@/memoryList?charset=utf8")
	Util.CheckError(openError)
	defer memoryDB.Close() //函数末关闭数据库
	// 设置最大打开的连接数，默认值为0表示不限制。
	memoryDB.SetMaxOpenConns(100)
	// 用于设置闲置的连接数。
	memoryDB.SetMaxIdleConns(50)
	memoryDB.Ping()
	//初始化数据库
	_, _ = memoryDB.Exec("CREATE TABLE IF NOT EXISTS memoList(memo TEXT)")
	_, _ = memoryDB.Exec("CREATE TABLE IF NOT EXISTS lastSetName(name TEXT)")

	var staticMemoList []string
	var lastSetName string
	// 查询多条数据
	rows, queryError := memoryDB.Query("SELECT memo FROM memoList")
	Util.CheckError(queryError)
	// 对多条数据进行遍历
	var memo string
	for rows.Next() {
		scanError := rows.Scan(&memo)
		Util.CheckError(scanError)
		//emoji表情解码
		fullTimeString := time.Now().String()
		timeString := fullTimeString[:19]
		fmt.Printf("TIME:%s --> 解码前数据:%s\n", timeString, memo)
		memo = Util.UnicodeEmojiDecode(memo)
		fullTimeString = time.Now().String()
		timeString = fullTimeString[:19]
		fmt.Printf("TIME:%s --> 解码后数据:%s\n", timeString, memo)

		staticMemoList = append(staticMemoList, memo)
	}

	nameRows, queryError := memoryDB.Query("SELECT name FROM lastSetName")
	Util.CheckError(queryError)
	var name string
	for nameRows.Next() {
		scanError := nameRows.Scan(&name)
		Util.CheckError(scanError)
		lastSetName = Util.UnicodeEmojiDecode(name)
	}

	nullMemoList := []string{}
	var currentMemoList []string
	if len(staticMemoList) <= 0 {
		currentMemoList = nullMemoList
	} else {
		currentMemoList = staticMemoList
	}

	fullTimeString := time.Now().String()
	timeString := fullTimeString[:19]
	fmt.Printf("TIME:%s --> 购物车最后上报人:%s\t | 响应的列表:%s\n", timeString, lastSetName, currentMemoList)
	ctx.JSON(iris.Map{"name": lastSetName, "memoList": currentMemoList})
}

func SetMemoListHadnler(ctx context.Context) {
	c := &MemoryModel.Memory{}
	if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
	} else {
		fullTimeString := time.Now().String()
		timeString := fullTimeString[:19]
		fmt.Printf("TIME:%s --> 购物车上报人:%s\t | 上报的列表:%#v\n", timeString, c.Name, c.MemoList)
		staticMemoList := c.MemoList
		lastSetName := c.Name
		//处理db
		db, openError := sql.Open("mysql", "root:simon.1314@/memoryList?charset=utf8") //本地数据库用户名root,密码:simon.1314
		Util.CheckError(openError)
		defer db.Close() //函数末关闭数据库

		// 设置最大打开的连接数，默认值为0表示不限制。
		db.SetMaxOpenConns(100)
		// 用于设置闲置的连接数。
		db.SetMaxIdleConns(50)
		db.Ping()
		_, execError1 := db.Exec("DROP TABLE IF EXISTS memoList")
		Util.CheckError(execError1)
		_, execError2 := db.Exec("CREATE TABLE memoList(memo TEXT)")
		Util.CheckError(execError2)

		//这边事物内批量数据插入
		tx, dbError := db.Begin()
		Util.CheckError(dbError)
		stmt, prepareError := tx.Prepare("INSERT memoList SET memo=?")
		Util.CheckError(prepareError)
		for _, value := range staticMemoList {
			//emoji表情转码
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
		_, execError3 := db.Exec("DROP TABLE IF EXISTS lastSetName")
		Util.CheckError(execError3)
		_, execError4 := db.Exec("CREATE TABLE lastSetName(name TEXT)")
		Util.CheckError(execError4)
		updateStmt, updateError := db.Prepare("INSERT INTO lastSetName(name) VALUES(?)")
		Util.CheckError(updateError)
		_, execUpdateErr := updateStmt.Exec(Util.UnicodeEmojiCode(lastSetName))
		Util.CheckError(execUpdateErr)

		ctx.JSON(iris.Map{"name": lastSetName, "memoList": staticMemoList})
	}
}
