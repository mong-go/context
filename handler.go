package gomon

import "net/http"
import "github.com/gorilla/context"
import "gopkg.in/mgo.v2"

// Handler provides a middlware handler to clone mgo sessions and set to context
// Alice constructor handler signature is returned
func Handler(ses *mgo.Session, name string, key ...interface{}) func(http.Handler) http.Handler {
	k := getkey(key...)
	if k == nil {
		k = name
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			cl := ses.Clone()
			defer cl.Close()
			db := cl.DB(name)
			context.Set(req, k, &db)

			h.ServeHTTP(w, req)
		})
	}
}

func getkey(key ...interface{}) interface{} {
	if len(key) > 0 {
		return key[0]
	}

	return nil
}
