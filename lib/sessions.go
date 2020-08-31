package lib

import "github.com/gorilla/sessions"

var store = sessions.NewCookieStore([]byte("secret_key"))

func GetStore() *sessions.CookieStore {
	return store
}
