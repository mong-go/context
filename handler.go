package context

import (
	"net/http"

	. "github.com/gorilla/context"
	"gopkg.in/mgo.v2"
)

type constructor func(http.Handler) http.Handler

// Handler provides a middlware handler to clone mgo sessions and set to context
func Handler(s *mgo.Session, name string, key ...interface{}) constructor {
	// get key or use the name of the db
	k := getkey(key...)
	if k == nil {
		k = name
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			c := s.Clone()
			defer c.Close()
			Set(req, k, c.DB(name))

			h.ServeHTTP(w, req)
		})
	}
}

// getkey returns the first key in the argument of keys
func getkey(key ...interface{}) interface{} {
	if len(key) > 0 {
		return key[0]
	}

	return nil
}
