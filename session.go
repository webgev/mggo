package mggo

import (
	"crypto/rand"
	"fmt"
	"time"
)

func init() {
	InitCallback(func() {
		models := []interface{}{(*SessionStorage)(nil)}
		CreateTable(models)
	})
}

// SessionStorage - session storage struct
type SessionStorage struct {
	Sid  string
	Date time.Time
}

// Set - set session sid
func (s SessionStorage) Set() {
	SQL().Insert(s)
}

// Delete - delete session sid
func (s SessionStorage) Delete() {
	SQL().Delete(s)
}

// GenerateSid - generate sid
func GenerateSid(id int) string {
	b := make([]byte, 16)
	rand.Read(b)
	sid := fmt.Sprintf("%d-%X-%X-%X-%X-%X", id, b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return sid
}
