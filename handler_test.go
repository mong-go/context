package gomon

import (
	"github.com/gorilla/context"
	"github.com/nowk/assert"
	"gopkg.in/mgo.v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mongod(t *testing.T, req *http.Request, fn func(s *mgo.Session)) {
	s, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		t.Fatalf("error: mgo dial: %s\n", err)
	}
	defer s.Close()

	fn(s)
}

var h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello World!"))
})

func TestHandlerContextAsDatabaseName(t *testing.T) {
	w, req := httptest.NewRecorder(), &http.Request{}

	mongod(t, req, func(s *mgo.Session) {
		Handler(s, "gomon_test")(h).ServeHTTP(w, req)
		se := context.Get(req, "gomon_test")
		assert.NotNil(t, se)
		assert.Equal(t, "Hello World!", w.Body.String())
	})
}

func TestHandlerSetCustomContextKey(t *testing.T) {
	w, req := httptest.NewRecorder(), &http.Request{}

	mongod(t, req, func(s *mgo.Session) {
		Handler(s, "gomon_test", "db")(h).ServeHTTP(w, req)
		se := context.Get(req, "db")
		assert.NotNil(t, se)
		assert.Equal(t, "Hello World!", w.Body.String())
	})
}
