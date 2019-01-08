package main

import(
	"fmt"
	
	"github.com/kataras/iris"
	// "github.com/kataras/iris/context"

	"sync"

"time"
)

//以下为初学测试
func newHello() {
	fmt.Println("hello world!")
}

var m *sync.RWMutex

func main(){
	// hello()
	m = new(sync.RWMutex)



//写的时候啥都不能干

go write(1)

go read(2)

go write(3)


time.Sleep(4 * time.Second)
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






func read(i int) {

println(i, "read start")



m.RLock()

println(i, "reading")

time.Sleep(1 * time.Second)

m.RUnlock()



println(i, "read end")
}

func write(i int) {

println(i, "write start")



m.Lock()

println(i, "writing")

time.Sleep(1 * time.Second)


m.Unlock()

println(i, "write end")
}
