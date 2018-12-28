package main

import (
	"fmt"
	"time"
	"Utils/util"

    "github.com/kataras/iris"
	"github.com/kataras/iris/context"

	"database/sql"
	_"github.com/Go-SQL-Driver/MySQL"
)

type Memory struct {
    Name  string
	MemoList []string
}
var staticMemoList []string 
var lastSetName string

func main() {
	memoryDB, openError := sql.Open("mysql", "root:simon.1314@/memoryList?charset=utf8")
	util.CheckError(openError)
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

	fullTimeString := time.Now().String()
	timeString := fullTimeString[:19]

	nameRows, queryError := memoryDB.Query("SELECT name FROM lastSetName")
	util.CheckError(queryError)
	var name string
	for nameRows.Next() {
		scanError := nameRows.Scan(&name)
		util.CheckError(scanError)
		lastSetName = util.UnicodeEmojiDecode(name);
	}

	//开启服务器
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
	fmt.Printf("TIME:%s --> 最后上报人:%s\t | 响应的列表:%s\n", timeString,lastSetName,currentMemoList)
	ctx.JSON(iris.Map{"name": lastSetName , "memoList" : currentMemoList})
}

func setMemoListHadnler(ctx context.Context) {
    c := &Memory{}
    if err := ctx.ReadJSON(c); err != nil {
		panic(err.Error())
    } else {
		fullTimeString := time.Now().String()
		timeString := fullTimeString[:19]
        fmt.Printf("TIME:%s --> 上报人:%s\t | 上报的列表:%#v\n",timeString,c.Name, c.MemoList)
		staticMemoList = c.MemoList
		lastSetName = c.Name
		//db
		db, openError := sql.Open("mysql", "root:simon.1314@/memoryList?charset=utf8")//本地数据库用户名root,密码:simon.1314
		util.CheckError(openError)
		_, execError1 := db.Exec("DROP TABLE IF EXISTS memoList")
		util.CheckError(execError1)
		_, execError2 := db.Exec("CREATE TABLE memoList(memo TEXT)")
		util.CheckError(execError2)
	
		// 批量数据插入
		tx, dbError := db.Begin()
		util.CheckError(dbError)
		stmt, prepareError := tx.Prepare("INSERT memolist SET memo=?")
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

func postForMemoListApp(){
	app := iris.New()
	//获取备忘录列表
	app.Post("/getMemoList", getMemoListHadnler)
	//全量更新备忘录列表
	app.Post("/setMemoList", setMemoListHadnler)
    app.Run(iris.Addr(":8080"))
}

