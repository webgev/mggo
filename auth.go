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
func (s SAP) Authenticate(login, password string) bool {
	if login != "" && password != "" {
		id := User{}.Identity(login, password)
		if id == 0 {
			return false
		}
		expiration := time.Now().Add(365 * 24 * time.Hour)
		sid := GenerateSid(id)
		cookie1 := &http.Cookie{Name: "sid", Value: sid, HttpOnly: true, Expires: expiration, Path: "/"}
		SetCookie(cookie1)
		session := SessionStorage{Sid: sid}
		session.Set()
		EventPublish("SAP.Auth", EventTypeServer, nil, id)
		return true
	}
	return false
}

// IsAuth is check user authenticate
func (s SAP) IsAuth() bool {
	return GetCookie("sid") != ""
}

// Exit is exit
func (s SAP) Exit() {
	expiration := time.Now().Add(-300)
	sid := GetCookie("sid")
	cookie1 := &http.Cookie{Name: "sid", Value: sid, HttpOnly: true, Expires: expiration, Path: "/"}
	SetCookie(cookie1)
	session := SessionStorage{Sid: sid}
	session.Delete()
}

// SessionUserID - get userid for cookie sid
func (s SAP) SessionUserID() int {
	sid := GetCookie("sid")
	if sid == "" {
		return 0
	}
	arr := strings.Split(sid, "-")
	id, _ := strconv.Atoi(arr[0])
	return id
}

// SessionID - get userid for cookie sid
func (s SAP) SessionID() string {
	return GetCookie("sid")
}
