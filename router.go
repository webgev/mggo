package mggo

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

var customHandles map[string]func(http.ResponseWriter, *http.Request)

// ParamsMethod is for api method
// Method is method name. Eq News.Read
// Params is map params for method.
// Parameters are passed to the structure of the method through the MapToString function.
type ParamsMethod struct {
	Params map[string]interface{}
	Method string
}
type resultMethod struct {
	Result interface{}
}

// RouterHooks is hooks api and view request
type RouterHooks interface {
	Before(*Router, http.ResponseWriter, *http.Request)
	After(*Router, http.ResponseWriter, *http.Request)
}

// Router web app
// - Host default localhost:9000
// - ViewData template view
// - Menu menu for template["Menu"]
type Router struct {
	ViewData    ViewData
	Menu        Menu
	RouterHooks RouterHooks
}

// run http
func (r *Router) run() {
	r.defaultParams()
	serverConfig, err := config.GetSection("server")
	if err != nil {
		panic(err)
	}
	host, err := serverConfig.GetKey("http_host")
	if err != nil {
		panic(err)
	}
	if add, err := serverConfig.GetKey("api_address"); err == nil {
		http.HandleFunc(add.String(), r.api)
	}
	if add, err := serverConfig.GetKey("view_address"); err == nil {
		http.HandleFunc(add.String(), r.view)
	}
	if add, err := serverConfig.GetKey("socket_address"); err == nil {
		http.HandleFunc(add.String(), func(w http.ResponseWriter, req *http.Request) {
			ctx := newBaseContext(w, req, nil, req.URL.Query(), User{})
			socketConnect(ctx, r.getUserInfo(ctx).ID)
		})
	}
	if static, err := serverConfig.GetKey("static_address"); err == nil {
		http.Handle(static.String(), http.StripPrefix(static.String(), http.FileServer(http.Dir("."+static.String()))))
	}
	if customHandles != nil {
		for path, handler := range customHandles {
			http.HandleFunc(path, handler)
		}
	}
	go runRPC(serverConfig.Key("rpc_address").String())
	err = http.ListenAndServe(host.String(), nil)
	if err != nil {
		panic(err)
	}
}

// HandleFunc is added handle in router
func (r *Router) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	if customHandles == nil {
		customHandles = map[string]func(http.ResponseWriter, *http.Request){path: handler}
	} else {
		customHandles[path] = handler
	}
}

func (r *Router) api(w http.ResponseWriter, req *http.Request) {
	ctx := newBaseContext(w, req, nil, req.URL.Query(), User{})
	startServer(w, req)
	defer endServer(ctx, r.ViewData)
	if req.Method != "POST" {
		panic(ErrorViewNotFound{})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var rec ParamsMethod
	r.parseBody(req, &rec)

	user := r.getUserInfo(ctx)
	if !CheckRight(rec.Method, user.Right, true) {
		panic(ErrorStatusForbidden{})
	}
	result := InvokeAPI(ctx, &rec)
	if result != nil {
		result = GetAPIResult(result)
	}

	r.hooksBefore(w, req)
	json.NewEncoder(w).Encode(resultMethod{Result: result})
	r.hooksAfter(w, req)
}

func (r *Router) view(w http.ResponseWriter, req *http.Request) {
	startServer(w, req)
	r.ViewData.Data = map[string]interface{}{}
	r.ViewData.Data["tempalteParser"] = tempalteParser{r.ViewData}

	path := strings.Split(req.URL.Path[1:], "/")
	if path[0] == "" {
		path[0] = "home"
	}
	// remove last empty item
	if path[len(path)-1] == "" {
		path = path[:len(path)-1]
	}
	rout := path[0]
	if len(path) == 1 {
		path = append(path, "index")
	}
	if path[1] == "" {
		path[1] = "index"
	}

	ctx := newBaseContext(w, req, path, req.URL.Query(), User{})
	defer endServer(ctx, r.ViewData)
	go SAP{}.Update(ctx)

	user := r.getUserInfo(ctx)
	if !CheckViewRight(strings.Title(rout), strings.Title(path[1]), user.Right, false) {
		panic(ErrorViewNotFound{})
	}
	InvokeView(ctx, rout, path[1], &r.ViewData)

	view := r.ViewData.DirView + r.ViewData.View
	t, err := template.ParseFiles(r.ViewData.DirView+r.ViewData.Template, view)
	if err != nil {
		panic(err)
	}
	r.ViewData.Data["UserInfo"] = user
	r.Menu.SetActivePage(rout)
	r.ViewData.Data["Menu"] = r.Menu.items
	r.ViewData.Data["Rights"] = Rights

	r.hooksBefore(w, req)
	t.Execute(w, r.ViewData.Data)
	r.hooksAfter(w, req)
}

func (r *Router) hooksBefore(w http.ResponseWriter, req *http.Request) {
	if r.RouterHooks != nil {
		r.RouterHooks.Before(r, w, req)
	}
}
func (r *Router) hooksAfter(w http.ResponseWriter, req *http.Request) {
	if r.RouterHooks != nil {
		r.RouterHooks.After(r, w, req)
	}
}

func (r *Router) parseBody(req *http.Request, rec *ParamsMethod) {
	if strings.HasPrefix(req.Header.Get("Content-Type"), "application/json") {
		_ = json.NewDecoder(req.Body).Decode(&rec)
	} else if strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data") {
		if err := req.ParseMultipartForm(1024 * 1024 * 8); err != nil {
			panic(ErrorMethodNotFound{})
		}
		rec.Method = req.FormValue("method")
		json.Unmarshal([]byte(req.FormValue("params")), &rec.Params)
		file, handler, err := req.FormFile("file")
		if err != nil {
			panic(ErrorMethodNotFound{})
		}
		rec.Params["File"] = File{File: file, FileHeader: *handler}
	}
}

func (r *Router) defaultParams() {
	if r.ViewData.Template == "" {
		panic("No set ViewData.Template")
	}
}

func (r *Router) getUserInfo(ctx *BaseContext) User {
	u := User{}
	return u.GetCurrentUserInfo(ctx)
}

type tempalteParser struct {
	viewData ViewData
}

func (p tempalteParser) Label(html string) template.HTML {
	return template.HTML(html)
}

func (p tempalteParser) File(html string, data interface{}) template.HTML {
	t, err := template.ParseFiles(p.viewData.DirView + html)
	if err != nil {
		return template.HTML("error tempalte " + html + " " + err.Error())
	}
	var tpl bytes.Buffer
	if err = t.Execute(&tpl, data); err != nil {
		return template.HTML("error tempalte " + html + " " + err.Error())
	}

	result := tpl.String()
	return template.HTML(result)
}
func (p tempalteParser) Slice(args ...interface{}) []interface{} {
	return args
}
func (p tempalteParser) Map() map[string]interface{} {
	return map[string]interface{}{}
}
func (p tempalteParser) SetMap(dict map[string]interface{}, key string, value interface{}) map[string]interface{} {
	dict[key] = value
	return dict
}
func (p tempalteParser) StructToMap(i interface{}) map[string]interface{} {
	return StructToMap(i)
}
