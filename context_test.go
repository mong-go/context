package context

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

func isContexted(t *testing.T, key interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c, err := Get(req, key)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, c)
		assert.TypeOf(t, "*mgo.Database", c)

		w.Write([]byte("Hello World!"))
	})
}

func TestHandlerContextAsDatabaseName(t *testing.T) {
	s, teardown := Setup(t)
	defer teardown()

	w := httptest.NewRecorder()
	req := &http.Request{}
	Handler(s, "gomon_test")(isContexted(t, "gomon_test")).ServeHTTP(w, req)
	assert.Equal(t, "Hello World!", w.Body.String())
}

func TestHandlerSetCustomContextKey(t *testing.T) {
	s, teardown := Setup(t)
	defer teardown()

	w := httptest.NewRecorder()
	req := &http.Request{}
	Handler(s, "gomon_test", "db")(isContexted(t, "db")).ServeHTTP(w, req)
	assert.Equal(t, "Hello World!", w.Body.String())
}
