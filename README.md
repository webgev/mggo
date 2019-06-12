# MgGo - Web Framework

Controller - api and view methods

## Get Started

Enter Go's path (format varies based on OS):

	cd $GOPATH

Install Mggo Example:

	go get -u github.com/webgev/mggo-example
	go run ./github.com/webgev/mggo-example/mail.go

Open http://localhost:9000 in your browser and you should see "It works!"

## Example

```go
import "github.com/webgev/mggo"
func main() {
    temp := core.ViewData {
        DirView: "./view/",
        Template: "_template.html",
        Data: map[string]interface{}{},
    }
   
    rout := core.Router{
        ViewData: temp,
        Menu: getMenu(),
    }
    cfg, err := ini.Load("./config.ini")
    if err != nil {
        os.Exit(1)
    }
    core.Run(rout, cfg)
}
```

Controller
```go
package controller

import (
	"strconv"

	"github.com/webgev/mggo"
)
func init() {
	mggo.RegisterController("news", NewNews)

	mggo.AppendRight("News.Read", mggo.RRightGuest)

	mggo.AppendViewRight("News.Update", mggo.RRightEditor)
	mggo.InitCallback(func() {
		mggo.CreateTable([]interface{}{(*News)(nil)})
	})
}
func NewNews() *News {
	return &News{}
}

type News struct {
	ID   int
	Name string
}

func (c *News) Read(ctx *mggo.BaseContext) News {
	mggo.SQL().Select(c)
	return *c
}
func (v News) IndexView(ctx *mggo.BaseContext, data *mggo.ViewData) {
	data.View = "news/news.html"
	data.Data["Title"] = "News"
	data.Data["News"] = v.List(ctx)
}
```

## New controller 

```shell
go run ..\lib\new-controller.go -name=News
```

## Depends

- https://github.com/mitchellh/mapstructure
- https://github.com/go-pg/pg
- https://github.com/gorilla/websocket
- https://github.com/go-ini/ini
