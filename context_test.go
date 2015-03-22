package context

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/nowk/assert.v2"
)

func Setup(t *testing.T) (*mgo.Session, func()) {
	s, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		t.Fatal(err)
	}

	return s, func() {
		s.Close()
	}
}

var h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello World!"))
})

func TestHandlerContextAsDatabaseName(t *testing.T) {
	s, teardown := Setup(t)
	defer teardown()

	w := httptest.NewRecorder()
	req := &http.Request{}
	Handler(s, "gomon_test")(h).ServeHTTP(w, req)

	c := context.Get(req, "gomon_test")
	assert.NotNil(t, c)
	assert.TypeOf(t, "*mgo.Database", c)
	assert.Equal(t, "Hello World!", w.Body.String())
}

func TestHandlerSetCustomContextKey(t *testing.T) {
	s, teardown := Setup(t)
	defer teardown()

	w := httptest.NewRecorder()
	req := &http.Request{}
	Handler(s, "gomon_test", "db")(h).ServeHTTP(w, req)

	c := context.Get(req, "db")
	assert.NotNil(t, c)
	assert.TypeOf(t, "*mgo.Database", c)
	assert.Equal(t, "Hello World!", w.Body.String())
}
