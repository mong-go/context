package context

import (
	"errors"
	"net/http"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	mgourl "gopkg.in/mong-go/url.v1"
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

			context.Set(req, k, c.DB(name))

			h.ServeHTTP(w, req)
		})
	}
}

// Parse parses a mongo url string, dials and returns the handler
func Parse(urlStr string, key ...interface{}) (constructor, error) {
	u, err := mgourl.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	s, err := mgo.Dial(u.ShortString())
	if err != nil {
		return nil, err
	}

	k := getkey(key...)
	if k == nil {
		k = u.Database()
	}

	return Handler(s, u.Database(), k), nil
}

// getkey returns the first key in the argument of keys
func getkey(key ...interface{}) interface{} {
	if len(key) > 0 {
		return key[0]
	}

	return nil
}

var ErrInvalidContext = errors.New("context was not *mgo.Database")

// Get returns the db from context or an invalid context error
func Get(req *http.Request, key interface{}) (*mgo.Database, error) {
	db, ok := context.Get(req, key).(*mgo.Database)
	if !ok {
		return nil, ErrInvalidContext
	}

	return db, nil
}
