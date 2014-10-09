package gomon

import "log"
import "net/http"
import "os"
import "github.com/gorilla/context"
import "gopkg.in/mgo.v2"

var databasename string

func init() {
	databasename = os.Getenv("MONGO_DB_NAME")
	if databasename == "" {
		log.Fatalln("error: database name: undefined")
	}
}

// Handler provides a middlware handler to clone mgo sessions and set to context
// Alice constructor handler signature is returned
func Handler(mgos *mgo.Session,
	ckey interface{}) func(http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			m := mgos.Clone()
			defer m.Close() // close after going down the alice chain

			db := m.DB(databasename)
			context.Set(req, ckey, &db)

			h.ServeHTTP(w, req)
		})
	}
}
