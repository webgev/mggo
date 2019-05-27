package mggo

import "errors"

// Field is record field
type Field struct {
    Value interface{}
    Name  string
}

// Get value
func (f *Field) Get() interface{} {
    return f.Value
}

// Set value
func (f *Field) Set(value interface{}) {
    f.Value = value
}

// Record struct
type Record struct {
    Fields map[string]*Field
}

// NewRecord create and return new Record
func NewRecord() *Record {
    return &Record{Fields: map[string]*Field{}}
}

// Add is added field in record
func (r *Record) Add(name string, value interface{}) error {
    if _, ok := r.Fields[name]; ok {
        return errors.New("column is exists")
    }
    r.Fields[name] = &Field{Value: value, Name: name}
    return nil
}

// Get is get field value
func (r *Record) Get(name string) interface{} {
    return r.Fields[name].Get()
}

// Set is set field value
func (r *Record) Set(name string, value interface{}) {
    r.Fields[name].Set(value)
}

// Has is check has field
func (r *Record) Has(name string) bool {
    _, ok := r.Fields[name]
    return ok
}

// ToMap is convert to map
func (r *Record) ToMap() map[string]interface{} {
    result := map[string]interface{}{}
    for key, value := range r.Fields {
        result[key] = value.Value
    }
    return result
}

// StructToRecord is convert to struct to record
func StructToRecord(i interface{}) *Record {
    m := StructToMap(i)
    r := NewRecord()

    for key, value := range m {
        r.Add(key, value)
    }
    return r
}
