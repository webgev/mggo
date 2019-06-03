package mggo

import "net/http"

// SetCookie - set coolie in response
func SetCookie(c *http.Cookie) {
	http.SetCookie(respose, c)
}

// GetCookie - set coolie in request
func GetCookie(name string) string {
	var cookie, err = request.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}
