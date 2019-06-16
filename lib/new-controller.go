// Create name controller

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	controllerStr = `package controller

import (
    "github.com/webgev/mggo"
    "strconv"
)
func init() {
	mggo.RegisterController("$NAMELOWER", New$NAME)

    mggo.AppendRight("$NAME.Read", mggo.RRightGuest)
    mggo.AppendRight("$NAME.List", mggo.RRightGuest)
    mggo.AppendRight("$NAME.Update", mggo.RRightEditor)
    mggo.AppendRight("$NAME.Delete", mggo.RRightEditor)

    mggo.AppendViewRight("$NAME.Update", mggo.RRightEditor)
    mggo.InitCallback(func () {
        mggo.CreateTable( []interface{}{ (*$NAME)(nil) } )
    })
}

func New$NAME() *$NAME {
	return &$NAME{}
}

type $NAME struct {
    ID int
    Name string
    mggo.ListFilter ` + "`sql:\"-\" structtomap:\"-\" mapstructure:\",squash\"`" + `
}

func(c $NAME) Read(ctx *mggo.BaseContext) $NAME {
    mggo.SQL().Select(&c)
    return c
}

func (c *$NAME) List(ctx *mggo.BaseContext) ($NAMELOWERs []$NAME) {
    query := mggo.SQL().Model(&$NAMELOWERs)
    for key, value := range c.Filter {
        switch key {
        case "Name":
            query.Where("name = ?", value)
        }
    }
    c.ListFilter.Paging(query).Select()
    return
}

func (c $NAME) Update(ctx *mggo.BaseContext) int {
    if c.ID == 0 {
        mggo.SQL().Insert(&c)
    } else {
        mggo.SQL().Update(&c)
    }
    return c.ID
}

func (c $NAME) Delete(ctx *mggo.BaseContext) {
    if c.ID != 0 {
        mggo.SQL().Delete(&c)
    }
}

func(c $NAME) IndexView(ctx *mggo.BaseContext, data *mggo.ViewData){
    data.View = "$NAMELOWER/$NAMELOWER.html"
    data.Data["Title"] = "$NAME"
    data.Data["$NAMEs"] = c.List(ctx)
}
func(v $NAME) ReadView(ctx *mggo.BaseContext, data *mggo.ViewData){
    if len(ctx.Path) > 2 {
        if i, err := strconv.Atoi(ctx.Path[2]); err == nil {
            data.View = "$NAMELOWER/read.html"
            c := $NAME{ID: i}
            r := c.Read(ctx)
            if r.ID > 0 {
                data.Data["Title"] = r.Name
                data.Data["$NAME"] = r
                return
            }
        } 
    }
    panic(mggo.ErrorViewNotFound{})
}
func(v $NAME) UpdateView(ctx *mggo.BaseContext, data *mggo.ViewData){
    data.View = "$NAMELOWER/Update.html"
    if len(ctx.Path) > 2 {
        if i, err := strconv.Atoi(ctx.Path[2]); err == nil {
            c := $NAME{ID: i,}
            r := c.Read(ctx)
            if r.ID == 0 {
                panic(mggo.ErrorViewNotFound{})
            }
            data.Data["Title"] = r.Name
            data.Data["$NAME"] = r	
        } else {
            panic(mggo.ErrorViewNotFound{})
        }
    } else {
        data.Data["Title"] = "Ceate $NAME" 
        data.Data["$NAME"] = $NAME{}
    }
}
`
	viewListStr = `{{define "content"}}
<table class="table table-striped">
    <thead>
        <tr>
            <th>ID</th>
            <th>Name</th>
        </tr>
    </thead>
    <tbody>
        {{range .$NAMEs}}
        <tr>
            <td>{{.ID}}</td>
            <td>{{.Name}}</td>
        </tr>
        {{end}}
    </tbody>
</table>
{{end}}
`
	viewReadStr = `{{define "content"}}
<div>
    <p>ID <strong>{{.$NAME.ID}} </strong> </h2>
    <p>Name <strong>{{.$NAME.Name}} </strong> </h2>
</div>
{{end}}
`
	viewUpdateStr = `{{define "content"}}
<div id="$NAME-update">
    <input type="text" data-bind="value: Name" placeholder="Name" />
    <input type="submit" title="Go" data-bind="click: clickHandler"/>  
</div>
<script>
    function $NAMEUpdate() {
        var self = this;
        this.Name = ko.observable({{.$NAME.Name}});
        this.ID = ko.observable({{.$NAME.ID}});
        this.clickHandler = function() {
            api("$NAME.Update", ko.toJS(self)).then(res=> {alert("ok")});
        }
    };
    var model = new $NAMEUpdate();
    ko.applyBindings(model, document.getElementById("$NAME-update"));
</script>
{{end}}
`
)

func main() {
	name := flag.String("name", "", "controller name")
	//isSql := flag.Bool("sql", true, "is sql")
	isView := flag.Bool("view", true, "is view")
	flag.Parse()

	if *name != "" {
		cName := strings.Title(*name)
		cNameLower := strings.ToLower(*name)
		text := strings.ReplaceAll(controllerStr, "$NAMELOWER", cNameLower)
		text = strings.ReplaceAll(text, "$NAME", cName)
		file, err := os.Create("./controller/" + cNameLower + ".go")
		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		defer file.Close()
		file.WriteString(text)

		if *isView {
			newpath := filepath.Join(".", "view", cNameLower)
			os.MkdirAll(newpath, os.ModePerm)
			filevList, err := os.Create(newpath + "/" + cNameLower + ".html")
			filevRead, err2 := os.Create(newpath + "/read.html")
			filevUpdate, err3 := os.Create(newpath + "/update.html")
			if err != nil || err2 != nil || err3 != nil {
				fmt.Println("Unable to create file:", err)
				os.Exit(1)
			}
			defer filevList.Close()
			defer filevRead.Close()
			defer filevUpdate.Close()
			filevList.WriteString(strings.ReplaceAll(viewListStr, "$NAME", cName))
			filevRead.WriteString(strings.ReplaceAll(viewReadStr, "$NAME", cName))
			filevUpdate.WriteString(strings.ReplaceAll(viewUpdateStr, "$NAME", cName))
		}

		fmt.Println("Done.")
	}
}
