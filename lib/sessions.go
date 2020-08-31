package lib

import "github.com/gorilla/sessions"

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// return a pointer to the cookie store
func GetStore() *sessions.CookieStore {
	return store
}
