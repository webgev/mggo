package mggo

import (
	"crypto/rand"
	"fmt"
	"strconv"
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
	Sid    string
	Expire time.Duration
	Data   string
}

// Set - set session sid
func (s *SessionStorage) Set() {
	redisClient.Set(s.Sid, s.Data, 0)
	redisClient.Expire(s.Sid, s.Expire)
}

// Isset - get session sid
func (s *SessionStorage) Isset() (string, error) {
	if s.Sid == "" {
		return "", fmt.Errorf("not sid")
	}
	val, err := redisClient.Get(s.Sid).Result()
	return val, err
}

// Delete - delete session sid
func (s *SessionStorage) Delete() {
	redisClient.Del(s.Sid)
}

// GenerateSid - generate sid
func GenerateSid(id int) string {
	x := string("00000000" + strconv.FormatUint(uint64(id), 16))
	b := make([]byte, 16)
	rand.Read(b)
	sid := fmt.Sprintf("%s-%v-%X-%X-%X-%X-%X", x[len(x)-8:], time.Now().Unix(), b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return sid
}
