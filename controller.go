package mggo

import (
	"io"
	"mime/multipart"
	"os"
	"reflect"
	"strconv"

	"github.com/go-pg/pg/orm"
	"github.com/mitchellh/mapstructure"
)

// Controller interface
type Controller interface{}

// ListFilter is struct for controller list method
type ListFilter struct {
	Filter MapStringAny `sql:"-" structtomap:"-" mapstructure:"filter"`
	Nav    MapStringAny `sql:"-" structtomap:"-" mapstructure:"nav"`
}

// Paging is paginator
func (l ListFilter) Paging(query *orm.Query) *orm.Query {
	for key, value := range l.Nav {
		switch key {
		case "page":
			query.Offset(int(value.(float64)))
		case "limit":
			query.Limit(int(value.(float64)))
		}
	}
	return query
}

// ViewData - struct view data
// Template - general view template
// View - content view
// Data - template data
type ViewData struct {
	DirView  string
	Template string
	View     string
	Data     map[string]interface{}
}

// MapToStruct decode map to struct
func MapToStruct(params map[string]interface{}, contr interface{}) {
	mapstructure.Decode(params, &contr)
}

// StructToMap encode struct to map
func StructToMap(i interface{}) (values map[string]interface{}) {
	values = map[string]interface{}{}
	iVal := reflect.Indirect(reflect.ValueOf(&i)).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// You ca use tags here...
		// tag := typ.Field(i).Tag.Get("tagname")
		// Convert each type into a string for the url.Values string map
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		values[typ.Field(i).Name] = v
	}
	return
}

// MapStringAny - for filter
type MapStringAny map[string]interface{}

// File - structure file
type File struct {
	File       multipart.File
	FileHeader multipart.FileHeader
}

// Save - save file in path
func (f File) Save(path string) bool {
	file, err := os.OpenFile(path+f.FileHeader.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return false
	}
	defer file.Close()
	io.Copy(file, f.File)
	return true
}

// GetAPIResult - возвращает результатирующие данные.
// исключает из структуры поля, помеченные тэгом structtomap:"-"
func GetAPIResult(params interface{}) interface{} {
	tp := reflect.TypeOf(params)
	result := map[string]interface{}{}

	switch tp.Kind() {
	case reflect.Struct:
		val := reflect.Indirect(reflect.ValueOf(&params)).Elem()
		for i := 0; i < tp.NumField(); i++ {
			valueField := val.Field(i)
			typeField := val.Type().Field(i)
			tag := typeField.Tag
			v, ok := tag.Lookup("structtomap")
			if ok {
				if v == "-" {
					continue
				}
			}
			result[typeField.Name] = GetAPIResult(valueField.Interface())
		}
	case reflect.Slice, reflect.Array:
		resultArr := []interface{}{}
		items := reflect.ValueOf(params)
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			resultArr = append(resultArr, GetAPIResult(item.Interface()))
		}
		return resultArr
	case reflect.Map:
		val := reflect.ValueOf(params)
		for _, e := range val.MapKeys() {
			v := val.MapIndex(e)
			switch t := e.Interface().(type) {
			case string:
				result[string(t)] = GetAPIResult(v.Interface())
			}
		}
		return result
	default:
		return params
	}
	return result
}

// Invoke controller method
func Invoke(ctx *BaseContext, controller Controller, methodName string) (result interface{}) {
	c := reflect.ValueOf(controller)
	m := c
	if c.Type().Kind() == reflect.Ptr {
		m = reflect.Indirect(c)
	}
	name := m.Type().Name()
	objects := name + "." + methodName
	issetCache := Cache.isset(objects)

	LogInfo("Вызов метода ", objects)

	if issetCache {
		if v, ok := Cache.getMethod(ctx, objects, controller); ok {
			LogInfo("Cache get")
			LogInfo("Конец метода ", objects)
			return v
		}
	}

	method := c.MethodByName(methodName)
	if !method.IsValid() {
		panic(ErrorMethodNotFound{})
	}
	inputs := []reflect.Value{
		reflect.ValueOf(ctx),
	}
	res := method.Call(inputs)
	if len(res) > 0 {
		result = res[0].Interface()
	}
	if issetCache {
		LogInfo("Cache set")
		Cache.setMethod(ctx, objects, result, controller)
	}
	LogInfo("Конец метода ", objects)
	return
}
