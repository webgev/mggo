# Web Framework Application

Controller - api and view methods

## Get Started

Enter Go's path (format varies based on OS):

	cd $GOPATH

Install Mggo Example:

	go get -u github.com/webgev/mggo-example

Open http://localhost:9000 in your browser and you should see "It works!"

## Example

```go
func main() {
    temp := core.ViewData {
        DirView: "./view/",
        Template: "_template.html",
        Data: map[string]interface{}{},
    }
   
    rout := core.Router{
        GetController: getController,
        ViewData: temp,
        Menu: getMenu(),
    }
    cfg, err := ini.Load("./config.ini")
    if err != nil {
        os.Exit(1)
    }
    core.Run(rout, cfg)
}
func getController(controllerName string) interface{} {
	switch strings.ToLower(controllerName) {
	case "home":
		return &controller.Home{}
	}
	
	return nil
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