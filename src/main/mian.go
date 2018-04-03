package main
import "fmt"
import "github.com/kataras/iris"

func main() {
    // hello()
    myInt := 2
    fmt.Println("数据:",myInt)
    hello2()
}

func hello(){
    app := iris.New()
    app.RegisterView(iris.HTML("./views", ".html"))
    app.Get("/hello",func(ctx iris.Context){
        ctx.ViewData("message","hello world")
        ctx.View("hello.html")
    })
    app.Run(iris.Addr(":8080"))
}

func hello2(){
    app := iris.New()
    app.Get("/user/{id:long}",func(ctx iris.Context){
        ctx.JSON(iris.Map{"message" : "myUserId"})
    })
    app.Run(iris.Addr(":8080"))
}