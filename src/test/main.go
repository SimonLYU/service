package main

import(
	"fmt"
	
	"github.com/kataras/iris"
	// "github.com/kataras/iris/context"
)

//以下为初学测试
func newHello() {
	fmt.Println("hello world!")
}

func main(){
	hello()
}

func hello() {
	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html"))
	app.Get("/hello", func(ctx iris.Context) {
		ctx.ViewData("message", "hello world")
		ctx.View("hello.html")
	})
	app.Run(iris.Addr(":8081"))
}

func hello2() {
	app := iris.New()
	app.Get("/user/{id:long}", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "myUserId"})
	})
	app.Run(iris.Addr(":8080"))
}

func hello3(){
	app := iris.New()
    app.RegisterView(iris.HTML("./views", ".html"))

    app.Get("/", func(ctx iris.Context) {
        ctx.ViewData("message", "Hello world!")
        ctx.View("hello.html")
    })

    app.Run(iris.Addr(":8080"))
}
