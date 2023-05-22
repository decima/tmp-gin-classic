package security

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
	"math/big"
	"sigomid/core/http/utils/cookies"
)

func formatCookieName(domain string) string {
	return fmt.Sprintf("csrf:%v", domain)
}
func Csrf(c *gin.Context, domain string) template.HTML {
	session := sessions.Default(c)
	CsrfToken, _ := generateRandomString(16)
	session.Set(formatCookieName(domain), CsrfToken)
	session.Save()
	return template.HTML(fmt.Sprintf(`<input type="hidden" name="_csrf" value="%v">`, CsrfToken))
}

func IsCsrfValid(c *gin.Context, domain string) bool {
	token := extractCsrfToken(c)
	session := sessions.Default(c)
	storedToken := session.Get(formatCookieName(domain))
	if token != storedToken {
		return false
	}
	return true
}

func IsCsrfValidWithFlash(c *gin.Context, domain string) bool {
	if !IsCsrfValid(c, domain) {
		cookies.AddErrorFlash(c, "Invalid csrf token")
		return false
	}
	return true
}

func extractCsrfToken(c *gin.Context) string {
	r := c.Request

	if t := r.FormValue("_csrf"); len(t) > 0 {
		return t
	} else if t := r.URL.Query().Get("_csrf"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-CSRF-TOKEN"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-XSRF-TOKEN"); len(t) > 0 {
		return t
	}
	return ""
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
