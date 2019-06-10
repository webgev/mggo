package mggo

import "net/http"

// SetCookie - set coolikie in response
func SetCookie(ctx *BaseContext, c *http.Cookie) {
	http.SetCookie(ctx.Response, c)
}

// GetCookie - set coolie in request
func GetCookie(ctx *BaseContext, name string) string {
	var cookie, err = ctx.Request.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}
