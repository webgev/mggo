package mggo

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// SAP contoller for auth
type SAP struct{}

// Authenticate by login and password
func (s SAP) Authenticate(ctx *BaseContext, login, password string) bool {
	if login != "" && password != "" {
		user := User{}
		id := user.Identity(login, password)
		if id == 0 {
			return false
		}
		expiration := time.Now().Add(365 * 24 * time.Hour)
		sid := GenerateSid(id)
		cookie1 := &http.Cookie{Name: "sid", Value: sid, HttpOnly: true, Expires: expiration, Path: "/"}
		SetCookie(ctx, cookie1)
		session := SessionStorage{Sid: sid}
		session.Set()
		EventPublish("SAP.Auth", EventTypeServer, nil, id)
		return true
	}
	return false
}

// IsAuth is check user authenticate
func (s SAP) IsAuth(ctx *BaseContext) bool {
	return GetCookie(ctx, "sid") != ""
}

// Exit is exit
func (s SAP) Exit(ctx *BaseContext) {
	expiration := time.Now().Add(-300)
	sid := GetCookie(ctx, "sid")
	cookie1 := &http.Cookie{Name: "sid", Value: sid, HttpOnly: true, Expires: expiration, Path: "/"}
	SetCookie(ctx, cookie1)
	session := SessionStorage{Sid: sid}
	session.Delete()
}

// SessionUserID - get userid for cookie sid
func (s SAP) SessionUserID(ctx *BaseContext) int {
	sid := GetCookie(ctx, "sid")
	if sid == "" {
		return 0
	}
	arr := strings.Split(sid, "-")
	id, _ := strconv.Atoi(arr[0])
	return id
}

// SessionID - get userid for cookie sid
func (s SAP) SessionID(ctx *BaseContext) string {
	return GetCookie(ctx, "sid")
}
