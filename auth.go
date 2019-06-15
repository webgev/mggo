package mggo

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// SAP contoller for auth
type SAP struct{}

func (s SAP) setSid(ctx *BaseContext, sid string) {
	expiration := time.Now().Add(10 * 365 * 24 * time.Hour)
	cookie1 := &http.Cookie{Name: "sid", Value: sid, HttpOnly: true, Expires: expiration, Path: "/"}
	SetCookie(ctx, cookie1)
	session := SessionStorage{Sid: sid, Expire: 24 * 30 * time.Hour}
	session.Set()
}

// Authenticate by login and password
func (s SAP) Authenticate(ctx *BaseContext, login, password string) bool {
	if login != "" && password != "" {
		user := User{}
		id := user.Identity(login, password)
		if id == 0 {
			return false
		}
		sid := GenerateSid(id)
		s.setSid(ctx, sid)
		EventPublish("SAP.Auth", EventTypeServer, nil, id)
		return true
	}
	return false
}

// IsAuth is check user authenticate
func (s SAP) IsAuth(ctx *BaseContext) bool {
	session := SessionStorage{Sid: s.SessionID(ctx)}
	_, err := session.Isset()
	return err == nil
}

// Exit is exit
func (s SAP) Exit(ctx *BaseContext) {
	expiration := time.Now().Add(-300)
	sid := s.SessionID(ctx)
	cookie1 := &http.Cookie{Name: "sid", Value: sid, HttpOnly: true, Expires: expiration, Path: "/"}
	SetCookie(ctx, cookie1)
	session := SessionStorage{Sid: sid}
	session.Delete()
}

// SessionUserID - get userid for cookie sid
func (s SAP) SessionUserID(ctx *BaseContext) int {
	if s.IsAuth(ctx) == false {
		return 0
	}
	sid := s.SessionID(ctx)
	arr := strings.Split(sid, "-")
	id, _ := strconv.ParseUint(arr[0], 16, 64)

	return int(id)
}

// Update session
func (s SAP) Update(ctx *BaseContext) {
	sid := s.SessionID(ctx)
	if sid == "" {
		return
	}
	s.setSid(ctx, sid)
}

// SessionID - get userid for cookie sid
func (s SAP) SessionID(ctx *BaseContext) string {
	return GetCookie(ctx, "sid")
}
