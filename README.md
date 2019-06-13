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
package main

import (
	"github.com/webgev/mggo"
	_ "./controller" // контроллеры
)


func main() {
	// конфигурируем view шаблон
	// DirView - папка, в которой находятся все шаблоны контроллеров
	// Template - основной шаблон 
	// Data - параметры
	temp := mggo.ViewData{
		DirView:  "./view/",
		Template: "_template.html",
		Data:     map[string]interface{}{},
	}

	rout := mggo.Router{
		ViewData:    temp,
	}

	mggo.Run(rout, "./config.ini")
}
```

config.ini
```ini
[server]
http_host = localhost:9000
view_address = /
api_address = /api/
socket_address = /echo
static_address = /static/
rpc_address = localhost:9001

[database]
user     = postgres
password = postgres
database = Go
address  = 127.0.0.1:5432
network  = tcp

[smtp]
email = *****@gmail.com
password = *****
server = smtp.gmail.com
port = 587

[redis]
address = localhost:6379
```

## Controller
```go
package controller

import (
	"strconv"

	"github.com/webgev/mggo"
)
func init() {
	// регистрируем контроллер
	mggo.RegisterController("news", NewNews)

	// добавляем права на api методы
	mggo.AppendRight("News.List", mggo.RRightGuest)
	mggo.AppendRight("News.Read", mggo.RRightGuest)
	mggo.AppendRight("News.Update", mggo.RRightManager)
	mggo.AppendRight("News.Delete", mggo.RRightManager)

	// добавляем права на view 
	mggo.AppendViewRight("News.Update", mggo.RRightManager)

	// после инициализации конфигурациии
	mggo.InitCallback(func() {
		//создаем таблицу
		models := []interface{}{(*mggo.News)(nil)}
		mggo.CreateTable(models)
		// кэшируем метод по параметрам на 2 часа
		mggo.Cache.AddMethod("News.List", mggo.CacheTypeMethodParams, 60*60*2)
	})
}
// Конструктор 
func NewNews() *News {
	return &News{}
}
// Оснавная структура новостей
// mggo.ListFilter - вспомогательная структура для списочных методов
// тег sql:"-" говорит о том, что данное поле не добавляется в таблицу
// тег structtomap:"-" - данное поле не будет возращаться при api запросе
// тег mapstructure:"id" - название поля при api запосе 
// тег mapstructure:",squash" - обрабатывем это поле при api запросе 
type News struct {
	ID  int `mapstructure:"id"`
	Name sting `mapstructure:"name"`
	mggo.ListFilter `sql:"-" structtomap:"-" mapstructure:",squash"`
}
// News.Read
func (n *News) Read(ctx *mggo.BaseContext) News {
	if n.ID != 0 {
		return News{}
	}
	mggo.SQL().Select(n)
	return *n
}
// News.List
func (n *News) List(ctx *mggo.BaseContext) (newsList []News) {
	query := mggo.SQL().Model(&newsList)
	for key, value := range n.Filter {
		switch key {
		case "name":
			query.Where("name = ?", value)
		}
	}
	n.ListFilter.Paging(query).Select()
	return
}
// News.Update
func (n *News) Update(ctx *mggo.BaseContext) News {
	if n.ID == 0 {
		mggo.SQL().Insert(n)
	} else {
		mggo.SQL().Update(n)
	}
	return *n
}
// News.Delete
func (n *News) Delete(ctx *mggo.BaseContext) {
	if n.ID != 0 {
		mggo.SQL().Delete(n)
	}
}
// News.ReadByName
func (n *News) ReadByName(ctx *mggo.BaseContext) News {
	if n.Name != "" {
		mggo.SQL().Model(n).Where("name = ?", n.Name).Select()
	}
	return *n
}

// View
// адрес - /news/
func (n News) IndexView(ctx *mggo.BaseContext, data *mggo.ViewData) {
	data.View = "news/news.html"
	data.Data["Title"] = "News"
	data.Data["Newss"] = n.List(ctx)
}
// адрес - /news/read/[newsID]
func (n News) ReadView(ctx *mggo.BaseContext, data *mggo.ViewData) {
	if len(ctx.Path) > 2 {
		if i, err := strconv.Atoi(ctx.Path[2]); err == nil {
			data.View = "news/read.html"
			c := News{ID: i}
			r := c.Read(ctx)
			if r.ID > 0 {
				data.Data["Title"] = r.Name
				data.Data["News"] = r
				return
			}
		}
	}
	panic(mggo.ErrorViewNotFound{})
}
// адрес - /news/update/  - создание новой новости
// /news/update[newsID]/[newsID]  - редактирование новости
func (n News) UpdateView(ctx *mggo.BaseContext, data *mggo.ViewData) {
	data.View = "news/update.html"
	if len(ctx.Path) > 2 {
		if i, err := strconv.Atoi(ctx.Path[2]); err == nil {
			data.View = "news/update.html"
			c := News{ID: i}
			r := c.Read(ctx)
			if r.ID == 0 {
				panic(mggo.ErrorViewNotFound{})
			}
			data.Data["Title"] = r.Name
			data.Data["News"] = r
		} else {
			panic(mggo.ErrorViewNotFound{})
		}
	} else {
		data.Data["Title"] = "Ceate News"
		data.Data["News"] = News{}
	}
}
// адрес - /news/[newsName]/  - поиск новости по имени 
func (n News) View(ctx *mggo.BaseContext, data *mggo.ViewData) {
	if len(ctx.Path) > 1 {
		n.Name = ctx.Path[1]
		news := n.ReadByName(ctx)
		if news.ID != 0 {
			data.View = "news/read.html"
			data.Data["Title"] = news.Name
			data.Data["News"] = news
			return
		}
	}
	panic(mggo.ErrorViewNotFound{})
}
```

## View

Шаблон _template.html
```html
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/css/bootstrap.min.css" integrity="sha384-GJzZqFGwb1QTTN6wy59ffF1BuGJpLSa9DkKMp0DgiMDm4iYMj70gZWKYbI706tWS" crossorigin="anonymous" />
<script type="text/javascript" src="/static/knockout-3.5.0.js"></script>
<script type="text/javascript" src="/static/app.js"></script>
<script type="text/javascript" src="/static/main.js"></script>
<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js" integrity="sha384-ZMP7rVo3mIykV+2+9J3UJ46jBk0WLaUAdn689aCwoqbBJiSnjAK/l8WvCWPIPm49" crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js" integrity="sha384-ChfqqxuZUCnJSK3+MXmPNIyE6ZbWh2IMqE241rYiqJxyMiZ6OW/JmZQ5stwEULTy" crossorigin="anonymous"></script>
</head>
<body>
{{.tempalteParser.File "layout/navbar.html" .}}
<div class="container">
	{{.tempalteParser.File "layout/userPanel.html" .}} 
	<h1> {{.Title}} </h1>
	{{template "content" .}}
</div>
</body>
</html>
```

navbar.html
```html
<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
	<div class="container">
		<a class="navbar-brand" href="/">Webgev</a>
		<button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarResponsive" aria-controls="navbarResponsive" aria-expanded="false" aria-label="Toggle navigation">
		<span class="navbar-toggler-icon"></span>
		</button>
		<div class="collapse navbar-collapse" id="navbarResponsive">
			<ul class="navbar-nav ml-auto">
				{{range .Menu}}
				<li class="nav-item {{if .Active}} active {{end}}">
					<a class="nav-link" href="{{ .Href }}">{{ .Title }}</a>
				</li>
				{{end}}
			</ul>
		</div>
	</div>
</nav>
```

userPanel.html
```html
<div id="userPanel" class="float-right">
{{if .UserInfo.ID}}
	<div class="btn btn-primary" data-bind="click: Exit"> Exit </div>
	<script>
		var UserInfo = {{.UserInfo}}
	</script>
{{else }}
	<a class="btn btn-primary" href="/auth"> Auth </a>
	<a class="btn btn-secondary" href="/reg"> Reg </a>
{{end}}
</div>
<script>
var userPanel = {
	Exit: function () {
		api("Auth.Exit").then(()=> location.replace("/"))
	}
}
ko.applyBindings(userPanel, document.getElementById("userPanel"));
</script>
```

### View контроллера новостей News
news.html
```html
{{define "content"}}
<table class="table table-striped">
	<thead>
		<tr>
			<th>ID</th>
			<th>Name</th>
			<th></th>
		</tr>
	</thead>
	<tbody>
		{{range .Newss}}
		<tr>
			<td>{{.ID}}</td>
			<td>{{.Name}}</td>
	
			<td> 
				<span onclick="productdelete({{.ID}})">Удалить</span>
			</td>
		</tr>
		{{end}}
	</tbody>
</table>
<script>
	function productdelete(id) {
		api("News.Delete", {"id": id})
	}
</script>
{{end}}
```

read.html
```html
{{define "content"}}
<div>
	<p>ID <strong>{{.News.ID}} </strong> </h2>
	<p>Name <strong>{{.News.Name}} </strong> </h2>
	
</div>
{{end}}
```

update.html
```html
{{define "content"}}
<div id="News-update">
	<input type="text" data-bind="value: name" placeholder="Name" />
	<input type="submit" title="Go" data-bind="click: clickHandler"/>  
</div>
<script>
	function NewsUpdate() {
		var self = this;
		this.name = ko.observable({{.News.Name}});
		this.id = ko.observable({{.News.ID}});
		this.clickHandler = function() {
			api("News.Update", ko.toJS(self)).then(res=> {alert("ok")});
	    }
	};
	var model = new NewsUpdate();
	ko.applyBindings(model, document.getElementById("News-update"));
</script>
{{end}}
```

## Fast Create New Controller 

```shell
go run ..\lib\new-controller.go -name=News
```

## Depends

- https://github.com/mitchellh/mapstructure
- https://github.com/go-pg/pg
- https://github.com/gorilla/websocket
- https://github.com/go-ini/ini
